package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"flamespot-api/src/config"
	"flamespot-api/src/model"
	"flamespot-api/src/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DashboardService interface {
	GetSummary() (model.DashboardSummary, error)
	ExportDashboard() ([]byte, error)
}

type dashboardService struct {
	db *mongo.Database
}

func NewDashboardService(db *mongo.Database) DashboardService {
	return &dashboardService{db: db}
}

func (s *dashboardService) GetSummary() (model.DashboardSummary, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	catCollection := s.db.Collection(config.CollectionCategories)
	catCursor, err := catCollection.Find(ctx, bson.M{})
	if err != nil {
		return model.DashboardSummary{}, err
	}
	defer catCursor.Close(ctx)

	var categories []model.Category
	if err = catCursor.All(ctx, &categories); err != nil {
		return model.DashboardSummary{}, err
	}

	categoryMap := make(map[string]string)
	categorySalesMap := make(map[string]int)

	for _, c := range categories {
		categoryMap[c.ID.Hex()] = c.Name
		categorySalesMap[c.ID.Hex()] = 0
	}

	currentTime := time.Now()
	currentPeriodStart := currentTime.AddDate(0, 0, -6).Truncate(24 * time.Hour).UnixMilli()
	previousPeriodStart := currentTime.AddDate(0, 0, -13).Truncate(24 * time.Hour).UnixMilli()

	orderCollection := s.db.Collection(config.CollectionOrders)
	orderCursor, err := orderCollection.Find(ctx, bson.M{
		"orderedAt": bson.M{"$gte": previousPeriodStart},
	})
	if err != nil {
		return model.DashboardSummary{}, err
	}
	defer orderCursor.Close(ctx)

	var orders []model.Order
	if err = orderCursor.All(ctx, &orders); err != nil {
		return model.DashboardSummary{}, err
	}

	var currRevenue, prevRevenue float64
	var currOrders, prevOrders int
	var currDineIn, prevDineIn int
	var currTakeaway, prevTakeaway int
	var currItems, prevItems int
	var currTurnSum, prevTurnSum int64
	var currTurnCount, prevTurnCount int64

	revenueByDay := make(map[string]float64)

	for _, order := range orders {
		if order.Status != model.StatusDeleted {
			if order.OrderedAt != nil {
				orderedTime := *order.OrderedAt
				isCurrent := orderedTime >= currentPeriodStart

				if isCurrent {
					currRevenue += order.SubTotal
					currOrders++

					switch order.Type {
					case model.TypeDineIn:
						currDineIn++
					case model.TypeTakeaway:
						currTakeaway++
					}

					for _, item := range order.Items {
						currItems += item.Quantity
						if _, exists := categorySalesMap[item.MenuItem.CategoryID]; exists {
							categorySalesMap[item.MenuItem.CategoryID] += item.Quantity
						}
					}
					orderTimeObj := time.UnixMilli(orderedTime)
					dayStr := orderTimeObj.Format("Mon")
					revenueByDay[dayStr] += order.SubTotal

					if order.BilledAt != nil && *order.BilledAt >= orderedTime {
						currTurnSum += (*order.BilledAt - orderedTime)
						currTurnCount++
					}
				} else {
					prevRevenue += order.SubTotal
					prevOrders++

					switch order.Type {
					case model.TypeDineIn:
						prevDineIn++
					case model.TypeTakeaway:
						prevTakeaway++
					}

					for _, item := range order.Items {
						prevItems += item.Quantity
					}
					if order.BilledAt != nil && *order.BilledAt >= orderedTime {
						prevTurnSum += (*order.BilledAt - orderedTime)
						prevTurnCount++
					}
				}
			}
		}
	}

	currAvgTicket, prevAvgTicket := 0.0, 0.0
	if currOrders > 0 {
		currAvgTicket = currRevenue / float64(currOrders)
	}
	if prevOrders > 0 {
		prevAvgTicket = prevRevenue / float64(prevOrders)
	}

	currAvgItems, prevAvgItems := 0.0, 0.0
	if currOrders > 0 {
		currAvgItems = float64(currItems) / float64(currOrders)
	}
	if prevOrders > 0 {
		prevAvgItems = float64(prevItems) / float64(prevOrders)
	}

	currTurn, prevTurn := 0, 0
	if currTurnCount > 0 {
		currTurn = int((currTurnSum / currTurnCount) / 60000)
	}
	if prevTurnCount > 0 {
		prevTurn = int((prevTurnSum / prevTurnCount) / 60000)
	}

	calcTrend := func(curr, prev float64, higherIsBetter bool) (string, bool) {
		if prev == 0 {
			if curr > 0 {
				return "+100.0%", higherIsBetter
			} else if curr < 0 {
				return "-100.0%", !higherIsBetter
			}
			return "0.0%", true
		}
		percent := ((curr - prev) / prev) * 100
		isPos := percent >= 0
		if !higherIsBetter {
			isPos = percent <= 0
		}
		return fmt.Sprintf("%+.1f%%", percent), isPos
	}

	revTrend, revPos := calcTrend(currRevenue, prevRevenue, true)
	ordTrend, ordPos := calcTrend(float64(currOrders), float64(prevOrders), true)
	tktTrend, tktPos := calcTrend(currAvgTicket, prevAvgTicket, true)
	dinTrend, dinPos := calcTrend(float64(currDineIn), float64(prevDineIn), true)
	takTrend, takPos := calcTrend(float64(currTakeaway), float64(prevTakeaway), true)
	trnTrend, trnPos := calcTrend(float64(currTurn), float64(prevTurn), false)
	itmTrend, itmPos := calcTrend(currAvgItems, prevAvgItems, true)
	cusTrend, cusPos := calcTrend(float64(currOrders), float64(prevOrders), true)

	buildStat := func(label, value, trend string, isPos bool) model.DashboardStat {
		s := strings.ToLower(label)
		s = strings.ReplaceAll(s, "-", " ")
		s = strings.ReplaceAll(s, ".", "")
		id := strings.Join(strings.Fields(s), "_")

		return model.DashboardStat{
			ID:         id,
			Label:      label,
			Value:      value,
			Trend:      trend,
			IsPositive: isPos,
		}
	}

	stats := []model.DashboardStat{
		buildStat("Total Revenue", fmt.Sprintf("₹%.2f", currRevenue), revTrend, revPos),
		buildStat("Total Orders Placed", strconv.Itoa(currOrders), ordTrend, ordPos),
		buildStat("Avg. Ticket Value", fmt.Sprintf("₹%.2f", currAvgTicket), tktTrend, tktPos),
		buildStat("Dine-In Orders", strconv.Itoa(currDineIn), dinTrend, dinPos),
		buildStat("Takeaway Orders", strconv.Itoa(currTakeaway), takTrend, takPos),
		buildStat("Average Table Turn", fmt.Sprintf("%d min", currTurn), trnTrend, trnPos),
		buildStat("Items Per Order", fmt.Sprintf("%.1f", currAvgItems), itmTrend, itmPos),
		buildStat("Active Customers", strconv.Itoa(currOrders), cusTrend, cusPos),
	}

	var dynamicDays []string
	for i := 6; i >= 0; i-- {
		day := currentTime.AddDate(0, 0, -i).Format("Mon")
		dynamicDays = append(dynamicDays, day)
	}

	var revenueOverview []model.RevenueData
	for _, day := range dynamicDays {
		revenueOverview = append(revenueOverview, model.RevenueData{
			Day:   day,
			Value: revenueByDay[day],
		})
	}

	type kv struct {
		Key   string
		Value int
	}
	var sortedCategories []kv
	for k, v := range categorySalesMap {
		sortedCategories = append(sortedCategories, kv{k, v})
	}

	sort.SliceStable(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Value > sortedCategories[j].Value
	})

	topCategories := []model.CategoryData{}
	for i, kv := range sortedCategories {
		if i >= 4 {
			break
		}

		catName := "Unknown"
		if name, exists := categoryMap[kv.Key]; exists {
			catName = name
		}

		topCategories = append(topCategories, model.CategoryData{
			Name:  catName,
			Value: kv.Value,
		})
	}

	return model.DashboardSummary{
		Stats:           stats,
		RevenueOverview: revenueOverview,
		TopCategories:   topCategories,
	}, nil
}

func (s *dashboardService) ExportDashboard() ([]byte, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionOrders)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []model.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"Order ID", "Type", "Status", "Base Total", "Discount", "Paid (SubTotal)", "Item Count", "Date"}
	_ = writer.Write(headers)

	for _, order := range orders {
		dateStr := ""
		if order.OrderedAt != nil {
			dateStr = time.UnixMilli(*order.OrderedAt).Format(time.RFC3339)
		}

		row := []string{
			order.ID,
			string(order.Type),
			string(order.Status),
			fmt.Sprintf("%.2f", order.Total),
			fmt.Sprintf("%.2f", order.Discount),
			fmt.Sprintf("%.2f", order.SubTotal),
			strconv.Itoa(len(order.Items)),
			dateStr,
		}
		_ = writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), nil
}
