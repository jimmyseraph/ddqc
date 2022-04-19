package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jimmyseraph/ddmc/api"
)

const (
	THREAD_NUM = 2 // 虚拟线程数量
)

func main() {
	headers := api.FormHeaders()
	cartIndexQueryString := api.FormCartIndexQueryString()
	baseBody := api.FormCheckOrderBody()
	fmt.Println("========================= 1.开始调用获取购物车接口 =========================")
	var cardIndexResponse string
	var err error
	i := 0
	for {
		i++
		fmt.Printf("  第%d次尝试......", i)
		cardIndexResponse, err = api.CallCartIndex(headers, cartIndexQueryString)
		if err != nil {
			fmt.Printf("失败\n")
			fmt.Printf("  接口调用失败：%v", err)
			continue
		}
		var json_data interface{}
		json.Unmarshal([]byte(cardIndexResponse), &json_data)
		if !json_data.(map[string]interface{})["success"].(bool) {
			fmt.Printf("失败，原因：%v\n", json_data.(map[string]interface{})["msg"])
			continue
		}
		fmt.Printf("成功\n")
		productList := json_data.(map[string]interface{})["data"].(map[string]interface{})["new_order_product_list"].([]interface{})
		if len(productList) == 0 {
			fmt.Printf("  购物车已无有效商品，结束程序，请重新添加购物车\n")
			return
		}
		products := productList[0].(map[string]interface{})["products"].([]interface{})
		if len(products) == 0 {
			fmt.Printf("  购物车已无有效商品，结束程序，请重新添加购物车\n")
			return
		}
		fmt.Printf("  购物车中的有效商品为：\n")
		for index, item := range products {
			p := item.(map[string]interface{})
			fmt.Printf("    %d. 名称：《%s》， 单价：%s， 数量：%f\n", index+1, p["product_name"], p["price"], p["count"])
		}
		break
	}
	fmt.Println("========================================================================")

	var card_index_response_json_data interface{}
	json.Unmarshal([]byte(cardIndexResponse), &card_index_response_json_data)
	productsIn := card_index_response_json_data.(map[string]interface{})["data"].(map[string]interface{})["new_order_product_list"]
	productsStr := api.FormCheckOrderProducts(productsIn.([]interface{}))

	fmt.Println("========================= 2.开始调用检查订单接口 =========================")
	var checkOrderResponse string
	i = 0
	for {
		i++
		fmt.Printf("  第%d次尝试......", i)
		checkOrderResponse, err = api.CallCheckOrder(headers, baseBody, productsStr)
		if err != nil {
			fmt.Printf("失败\n")
			fmt.Printf("  接口调用失败：%v", err)
			continue
		}
		var json_data interface{}
		json.Unmarshal([]byte(checkOrderResponse), &json_data)
		if !json_data.(map[string]interface{})["success"].(bool) {
			fmt.Printf("失败，原因：%v\n", json_data.(map[string]interface{})["msg"])
			continue
		}
		fmt.Printf("成功\n")
		break
	}
	fmt.Println("========================================================================")

	fmt.Println("========================= 3.开始调用配送时间接口 =========================")
	productsStr_getMultiReserveTime := api.FormGetMultiReserveTimeProducts(productsIn.([]interface{}))
	// var getMultiReserveTimeResponse string
	var reserved_time_start, reserved_time_end uint32 = 0, 0
	i = 0
	for {
		i++
		fmt.Printf("  第%d次尝试......", i)
		getMultiReserveTimeResponse, err := api.CallGetMultiReserveTime(headers, baseBody, productsStr_getMultiReserveTime)
		if err != nil {
			fmt.Printf("失败\n")
			fmt.Printf("  接口调用失败：%v\n", err)
			continue
		}
		var json_data interface{}
		json.Unmarshal([]byte(getMultiReserveTimeResponse), &json_data)
		if !json_data.(map[string]interface{})["success"].(bool) {
			fmt.Printf("失败，原因：%v\n", json_data.(map[string]interface{})["msg"])
			continue
		}
		fmt.Printf("成功\n")
		data := json_data.(map[string]interface{})["data"].([]interface{})
		if len(data) == 0 {
			fmt.Printf("  未能获取配送时间\n")
			continue
		}
		time := data[0].(map[string]interface{})["time"].([]interface{})
		if len(time) == 0 {
			fmt.Printf("  未能获取配送时间\n")
			continue
		}
		// fmt.Printf("--> %s\n", getMultiReserveTimeResponse)
		reserved_time_start, reserved_time_end = api.GetAvailableReservedTime(time[0].(map[string]interface{})["times"].([]interface{}))
		if reserved_time_start == 0 {
			fmt.Printf("  当前运力紧张，今天各时段已约满，程序结束。\n")
			return
		}
		fmt.Printf("  存在有效配送时间！\n")
		break
	}
	fmt.Println("========================================================================")

	fmt.Println("========================= 4.开始调用创建新订单接口 =========================")
	var check_order_response_json_data interface{}
	json.Unmarshal([]byte(checkOrderResponse), &check_order_response_json_data)
	order := check_order_response_json_data.(map[string]interface{})["data"].(map[string]interface{})["order"]
	// 开始做协程
	var wg sync.WaitGroup
	wg.Add(THREAD_NUM)
	var stop bool = false
	for i := 0; i < THREAD_NUM; i++ {
		go func(index int) {
			defer wg.Done()
			thread_count := index + 1
			counter := 0
			for !stop {
				counter++
				fmt.Printf("  协程(%d)发起请求，第%d次尝试......", thread_count, counter)
				packageOrder := api.FormAddNewOrderPackageOrder(
					productsIn.([]interface{}),
					card_index_response_json_data.(map[string]interface{})["data"].(map[string]interface{})["parent_order_info"],
					order.(map[string]interface{})["total_money"].(string),
					order.(map[string]interface{})["freight_discount_money"].(string),
					order.(map[string]interface{})["freight_money"].(string),
					order.(map[string]interface{})["freight_real_money"].(string),
					reserved_time_start, reserved_time_end,
				)
				addNewOrderResponse, err := api.CallAddNewOrder(headers, baseBody, packageOrder)
				if err != nil {
					fmt.Printf("失败, 接口调用失败: %v\n", err)
					continue
				}
				var json_data interface{}
				json.Unmarshal([]byte(addNewOrderResponse), &json_data)
				if !json_data.(map[string]interface{})["success"].(bool) {
					if json_data.(map[string]interface{})["code"].(float64) != -3001 {
						fmt.Printf("失败，原因：%v\n", json_data.(map[string]interface{})["msg"])
					} else {
						fmt.Printf("失败，原因：%v\n", json_data.(map[string]interface{})["tips"].(map[string]interface{})["limitMsg"])
					}
					continue
				}
				fmt.Printf("成功，请及时去微信小程序完成付款\n")
				stop = true
			}

		}(i)
	}
	wg.Wait()
	fmt.Println("========================================================================")
}
