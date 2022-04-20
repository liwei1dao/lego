package influxdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	iclient "github.com/influxdata/influxdb-client-go/v2"
	"github.com/liwei1dao/lego/sys/influxdb"
)

func Test_sys_write(t *testing.T) {
	if sys, err := influxdb.NewSys(
		influxdb.SetAddr("http://172.20.27.145:8086"),
		influxdb.SetToken("1gvhucIwiarXQIgRYCU6Fvd64ZHDnNWcpXLR9753W_VeuC9YwG7eSSO20FDohtS6y4RLb0mgiPHD5DzGZRTApg=="),
	); err != nil {
		fmt.Printf("init influxdb err:%v", err)
		return
	} else {
		writeAPI := sys.WriteAPI("liwei1dao", "liwei1dao")
		p := iclient.NewPoint("stat",
			map[string]string{"unit": "temperature"},
			map[string]interface{}{"min": 30.0, "max": 35.0},
			time.Now())
		writeAPI.WritePoint(p)
		writeAPI.Flush()
		defer sys.Close()
	}
}
func Test_sys_write1(t *testing.T) {
	if sys, err := influxdb.NewSys(
		influxdb.SetAddr("http://172.20.27.145:8086"),
		influxdb.SetToken("1gvhucIwiarXQIgRYCU6Fvd64ZHDnNWcpXLR9753W_VeuC9YwG7eSSO20FDohtS6y4RLb0mgiPHD5DzGZRTApg=="),
	); err != nil {
		fmt.Printf("init influxdb err:%v", err)
		return
	} else {
		writeAPI := sys.WriteAPI("liwei1dao", "liwei1dao")
		p := iclient.NewPointWithMeasurement("stat").
			AddTag("unit", "temperature").
			AddField("状态", "正常").
			AddField("max", 35.0)
		writeAPI.WritePoint(p)
		writeAPI.Flush()
		defer sys.Close()
	}
}

func Test_sys_query(t *testing.T) {
	if sys, err := influxdb.NewSys(
		influxdb.SetAddr("http://172.20.27.145:8086"),
		influxdb.SetToken("1gvhucIwiarXQIgRYCU6Fvd64ZHDnNWcpXLR9753W_VeuC9YwG7eSSO20FDohtS6y4RLb0mgiPHD5DzGZRTApg=="),
	); err != nil {
		fmt.Printf("init influxdb err:%v", err)
		return
	} else {
		queryAPI := sys.QueryAPI("liwei1dao")
		query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -3h) |> filter(fn: (r) => r._measurement == \"stat\")", "liwei1dao")
		result1, err := queryAPI.QueryRaw(context.Background(), query, iclient.DefaultDialect())
		if err != nil {
			fmt.Printf("QueryAPI err:%v", err)
			return
		} else {
			fmt.Printf("result:%v", result1)
		}
		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			fmt.Printf("QueryAPI err:%v", err)
			return
		}
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		defer sys.Close()
	}
}
