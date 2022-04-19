package api

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jimmyseraph/sparkle/easy_http"
)

type Maicai_API struct {
	Url    string
	Method easy_http.Method
}

const (
	BASE_URL                    = "https://maicai.api.ddxq.mobi"
	CHECK_ORDER_PATH            = "/order/checkOrder"          // 检查订单
	GET_MULTI_RESERVE_TIME_PATH = "/order/getMultiReserveTime" // 获取配送时间
	CART_INDEX_PATH             = "/cart/index"                // 购物车清单
	ADD_NEW_ORDER_PATH          = "/order/addNewOrder"         // 下单
)

var CartIndexApi Maicai_API = Maicai_API{
	Url:    BASE_URL + CART_INDEX_PATH,
	Method: easy_http.GET,
}

var CheckOrderApi Maicai_API = Maicai_API{
	Url:    BASE_URL + CHECK_ORDER_PATH,
	Method: easy_http.POST,
}

var GetMultiReserveTimeApi Maicai_API = Maicai_API{
	Url:    BASE_URL + GET_MULTI_RESERVE_TIME_PATH,
	Method: easy_http.POST,
}

var AddNewOrderApi Maicai_API = Maicai_API{
	Url:    BASE_URL + ADD_NEW_ORDER_PATH,
	Method: easy_http.POST,
}

func CallCartIndex(headers map[string][]string, queryString string) (string, error) {
	handler := easy_http.NewGet(CartIndexApi.Url + "?" + queryString)
	handler.Headers = headers
	resp, err := handler.Execute()
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func CallCheckOrder(headers map[string][]string, body string, packages string) (string, error) {
	body = fmt.Sprintf("%s&address_id=&user_ticket_id=default&freight_ticket_id=default&is_use_point=0&is_use_balance=0&is_buy_vip=0&coupons_id=&is_buy_coupons=0&packages=%s&check_order_type=0&is_support_merge_payment=1&showData=true&showMsg=false",
		body, url.QueryEscape(packages),
	)
	// fmt.Println(body)
	handler := easy_http.NewPost(CheckOrderApi.Url, body)
	handler.Headers = headers
	handler.Headers["content-type"] = []string{"application/x-www-form-urlencoded"}
	resp, err := handler.Execute()
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func CallGetMultiReserveTime(headers map[string][]string, body string, packages string) (string, error) {
	body = fmt.Sprintf("%s&address_id=%s&group_config_id=&products=%s&isBridge=false",
		body, url.QueryEscape(ADDRESS_ID), url.QueryEscape(packages),
	)
	// fmt.Println(body)
	handler := easy_http.NewPost(GetMultiReserveTimeApi.Url, body)
	handler.Headers = headers
	handler.Headers["content-type"] = []string{"application/x-www-form-urlencoded"}
	resp, err := handler.Execute()
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func CallAddNewOrder(headers map[string][]string, body string, package_order string) (string, error) {
	body = fmt.Sprintf("%s&time=%d&package_order=%s&showData=true&showMsg=false&ab_config=%s",
		body, time.Now().UnixMilli(), url.QueryEscape(package_order), url.QueryEscape("{\"key_onion\":\"C\"}"),
	)
	handler := easy_http.NewPost(AddNewOrderApi.Url, body)
	for key, value := range headers {
		handler.Headers[key] = value
	}
	// handler.Headers = headers
	handler.Headers["content-type"] = []string{"application/x-www-form-urlencoded"}
	resp, err := handler.Execute()
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func FormHeaders() map[string][]string {
	headers := make(map[string][]string)
	headers["ddmc-city-number"] = []string{CITY_NUMBER}
	headers["ddmc-build-version"] = []string{APP_VERSION}
	headers["ddmc-device-id"] = []string{OPEN_ID}
	headers["ddmc-station-id"] = []string{STATION_ID}
	headers["ddmc-channel"] = []string{"applet"}
	headers["ddmc-app-client-id"] = []string{APP_CLIENT_ID}
	headers["cookie"] = []string{"DDXQSESSID=" + S_ID}
	headers["ddmc-longitude"] = []string{LONGITUDE}
	headers["ddmc-latitude"] = []string{LATITUDE}
	headers["ddmc-api-version"] = []string{API_VERSION}
	headers["ddmc-uid"] = []string{UID}
	headers["user-agent"] = []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x1800123f) NetType/WIFI Language/en"}
	return headers
}

func FormCartIndexQueryString() string {
	return fmt.Sprintf(`uid=%s&longitude=%s&latitude=%s&station_id=%s&city_number=%s&api_version=%s&app_version=%s&applet_source=&channel=applet&app_client_id=%s&sharer_uid=&s_id=%s&openid=%s&h5_source=&device_token=%s&is_load=1`,
		url.QueryEscape(UID), url.QueryEscape(LONGITUDE), url.QueryEscape(LATITUDE), url.QueryEscape(STATION_ID),
		url.QueryEscape(CITY_NUMBER), url.QueryEscape(API_VERSION), url.QueryEscape(APP_VERSION),
		url.QueryEscape(APP_CLIENT_ID), url.QueryEscape(S_ID), url.QueryEscape(OPEN_ID), url.QueryEscape(DEVICE_TOKEN))
}

func FormCheckOrderBody() string {
	return fmt.Sprintf(`uid=%s&longitude=%s&latitude=%s&station_id=%s&city_number=%s&api_version=%s&app_version=%s&applet_source=&channel=applet&app_client_id=%s&sharer_uid=&s_id=%s&openid=%s&h5_source=&device_token=%s`,
		url.QueryEscape(UID), url.QueryEscape(LONGITUDE), url.QueryEscape(LATITUDE), url.QueryEscape(STATION_ID),
		url.QueryEscape(CITY_NUMBER), url.QueryEscape(API_VERSION), url.QueryEscape(APP_VERSION),
		url.QueryEscape(APP_CLIENT_ID), url.QueryEscape(S_ID), url.QueryEscape(OPEN_ID), url.QueryEscape(DEVICE_TOKEN))
}

func FormCheckOrderProducts(productsIn []interface{}) string {
	products := make([]interface{}, len(productsIn))
	for index, item := range productsIn {
		ps := item.(map[string]interface{})["products"].([]interface{})
		products[index] = map[string]interface{}{
			"total_money":               item.(map[string]interface{})["total_money"],
			"total_origin_money":        item.(map[string]interface{})["total_origin_money"],
			"goods_real_money":          item.(map[string]interface{})["goods_real_money"],
			"total_count":               item.(map[string]interface{})["total_count"],
			"cart_count":                item.(map[string]interface{})["cart_count"],
			"is_presale":                item.(map[string]interface{})["is_presale"],
			"instant_rebate_money":      item.(map[string]interface{})["instant_rebate_money"],
			"used_balance_money":        item.(map[string]interface{})["used_balance_money"],
			"can_used_balance_money":    item.(map[string]interface{})["can_used_balance_money"],
			"used_point_num":            item.(map[string]interface{})["used_point_num"],
			"used_point_money":          item.(map[string]interface{})["used_point_money"],
			"can_used_point_num":        item.(map[string]interface{})["can_used_point_num"],
			"can_used_point_money":      item.(map[string]interface{})["can_used_point_money"],
			"is_share_station":          item.(map[string]interface{})["is_share_station"],
			"only_today_products":       item.(map[string]interface{})["only_today_products"],
			"only_tomorrow_products":    item.(map[string]interface{})["only_tomorrow_products"],
			"package_type":              item.(map[string]interface{})["package_type"],
			"package_id":                item.(map[string]interface{})["package_id"],
			"front_package_text":        item.(map[string]interface{})["front_package_text"],
			"front_package_type":        item.(map[string]interface{})["front_package_type"],
			"front_package_stock_color": item.(map[string]interface{})["front_package_stock_color"],
			"front_package_bg_color":    item.(map[string]interface{})["front_package_bg_color"],
			"reserved_time": map[string]interface{}{
				"reserved_time_start": nil,
				"reserved_time_end":   nil,
			},
			"products": make([]interface{}, len(ps)),
		}

		for i, p := range ps {
			products[index].(map[string]interface{})["products"].([]interface{})[i] = map[string]interface{}{
				"id":                   p.(map[string]interface{})["id"],
				"category_path":        p.(map[string]interface{})["category_path"],
				"count":                p.(map[string]interface{})["count"],
				"price":                p.(map[string]interface{})["price"],
				"total_money":          p.(map[string]interface{})["total_price"],
				"instant_rebate_money": p.(map[string]interface{})["instant_rebate_money"],
				"activity_id":          p.(map[string]interface{})["activity_id"],
				"conditions_num":       p.(map[string]interface{})["conditions_num"],
				"product_type":         p.(map[string]interface{})["product_type"],
				"sizes":                p.(map[string]interface{})["sizes"],
				"type":                 p.(map[string]interface{})["type"],
				"total_origin_money":   p.(map[string]interface{})["total_origin_money"],
				"price_type":           p.(map[string]interface{})["price_type"],
				"batch_type":           p.(map[string]interface{})["batch_type"],
				"sub_list":             p.(map[string]interface{})["sub_list"],
				"order_sort":           p.(map[string]interface{})["order_sort"],
				"origin_price":         p.(map[string]interface{})["origin_price"],
			}
		}
	}
	p, err := json.Marshal(products)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return ""
	}
	return string(p)
}

func FormGetMultiReserveTimeProducts(productsIn []interface{}) string {
	products := make([]interface{}, len(productsIn))
	for index, item := range productsIn {
		ps := item.(map[string]interface{})["products"].([]interface{})
		products[index] = ps
	}

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(products)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return ""
	}
	p := buffer.Bytes()
	return string(p)
}

func GetAvailableReservedTime(reserveTime []interface{}) (uint32, uint32) {
	var reserved_time_start, reserved_time_end uint32 = 0, 0
	for _, item := range reserveTime {
		if !item.(map[string]interface{})["fullFlag"].(bool) {
			reserved_time_start = uint32(item.(map[string]interface{})["start_timestamp"].(float64))
			reserved_time_end = uint32(item.(map[string]interface{})["end_timestamp"].(float64))
			break
		}
	}
	return reserved_time_start, reserved_time_end
}

func FormAddNewOrderPackageOrder(
	productsIn []interface{},
	parent_order_info interface{},
	price, freight_discount_money, freight_money, order_freight interface{},
	reserved_time_start, reserved_time_end uint32,
) string {
	// fmt.Printf("--> %v", parent_order_info)
	order := map[string]interface{}{
		"payment_order": map[string]interface{}{
			"reserved_time_start":    reserved_time_start,
			"reserved_time_end":      reserved_time_end,
			"price":                  price,
			"freight_discount_money": freight_discount_money,
			"freight_money":          freight_money,
			"order_freight":          order_freight,
			"parent_order_sign":      parent_order_info.(map[string]interface{})["parent_order_sign"],
			"product_type":           1,
			"address_id":             ADDRESS_ID,
			"form_id":                getFormId(),
			"receipt_without_sku":    nil,
			"pay_type":               6,
			"vip_money":              "",
			"vip_buy_user_ticket_id": "",
			"coupons_money":          "",
			"coupons_id":             "",
		},
		"packages": make([]interface{}, len(productsIn)),
	}
	for index, item := range productsIn {
		ps := item.(map[string]interface{})["products"].([]interface{})
		other := item.(map[string]interface{})
		order["packages"].([]interface{})[index] = map[string]interface{}{
			"products":                  ps,
			"total_money":               other["total_money"],
			"total_origin_money":        other["total_origin_money"],
			"goods_real_money":          other["goods_real_money"],
			"total_count":               other["total_count"],
			"cart_count":                other["cart_count"],
			"is_presale":                other["is_presale"],
			"instant_rebate_money":      other["instant_rebate_money"],
			"used_balance_money":        other["used_balance_money"],
			"can_used_balance_money":    other["can_used_balance_money"],
			"used_point_num":            other["used_point_num"],
			"used_point_money":          other["used_point_money"],
			"can_used_point_num":        other["can_used_point_num"],
			"can_used_point_money":      other["can_used_point_money"],
			"is_share_station":          other["is_share_station"],
			"only_today_products":       other["only_today_products"],
			"only_tomorrow_products":    other["only_tomorrow_products"],
			"package_type":              other["package_type"],
			"package_id":                other["package_id"],
			"front_package_text":        other["front_package_text"],
			"front_package_type":        other["front_package_type"],
			"front_package_stock_color": other["front_package_stock_color"],
			"front_package_bg_color":    other["front_package_bg_color"],
			"eta_trace_id":              "",
			"reserved_time_start":       reserved_time_start,
			"reserved_time_end":         reserved_time_end,
			"soon_arrival":              "",
			"first_selected_big_time":   1,
		}
	}
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(order)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return ""
	}
	p := buffer.Bytes()
	return string(p)
}

func getFormId() string {
	t := time.Now().UnixMilli()
	hash := crypto.MD5.New()
	hash.Write([]byte(strconv.FormatInt(t, 10)))
	b := hash.Sum(nil)
	return hex.EncodeToString(b)
}
