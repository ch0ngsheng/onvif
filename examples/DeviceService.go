package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	sdk "github.com/use-go/onvif/sdk/device"
	sdkMedia "github.com/use-go/onvif/sdk/media"
	sdkptz "github.com/use-go/onvif/sdk/ptz"
	"github.com/use-go/onvif/xsd/onvif"
)

const (
	xytLogin    = "admin"
	xytPassword = "xyt123456"
	xytAddr     = "xyt-camera.dtiot.com:46953"

	jtsqLogin    = "admin"
	jtsqPassword = "Zyf@2022"
	jtsqAddr     = "192.168.200.156"

	// mark 大华监控 老版本的onvif和登录密码不一样，且不支持修改
	zhtLogin    = "admin"
	zhgPassword = "admin"
	zhtAddr     = "10.11.202.210"
)

var (
	login    = zhtLogin
	password = zhgPassword
	addr     = zhtAddr
)

func main() {
	ctx := context.Background()

	//Getting an camera instance
	dev, err := goonvif.NewDevice(goonvif.DeviceParams{
		Xaddr:      addr,
		Username:   login,
		Password:   password,
		HttpClient: new(http.Client),
	})
	if err != nil {
		fmt.Printf("%+v", err)
		panic(err)
	}

	//Preparing commands
	systemDateAndTyme := device.GetSystemDateAndTime{}
	getCapabilities := device.GetCapabilities{Category: "All"}

	//Commands execution
	systemDateAndTymeResponse, err := sdk.Call_GetSystemDateAndTime(ctx, dev, systemDateAndTyme)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(systemDateAndTymeResponse)
	}
	getCapabilitiesResponse, err := sdk.Call_GetCapabilities(ctx, dev, getCapabilities)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(getCapabilitiesResponse)
	}

	p := media.GetProfiles{}
	profiles, err := sdkMedia.Call_GetProfiles(ctx, dev, p)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(profiles)
	}

	profile, err := sdkMedia.Call_GetProfile(ctx, dev, media.GetProfile{
		ProfileToken: "MediaProfile000",
	})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(profile)
	}

	uri, err := sdkMedia.Call_GetStreamUri(ctx, dev, media.GetStreamUri{
		ProfileToken: "MediaProfile000",
	})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(uri)
	}
	/*
		relativeMove, err := sdkptz.Call_RelativeMove(ctx, dev, ptz.RelativeMove{
			XMLName:      "",
			ProfileToken: "MediaProfile000",
			Translation: onvif.PTZVector{
				PanTilt: onvif.Vector2D{
					X:     0,
					Y:     0,
					Space: "",
				},
				Zoom: onvif.Vector1D{
					X:     -5,
					Space: "",
				},
			},
			Speed: onvif.PTZSpeed{},
		})
		if err != nil {
			log.Println(err)
		} else {
			log.Println(relativeMove)
		}
	*/

	moveResponse, err := sdkptz.Call_ContinuousMove(ctx, dev, ptz.ContinuousMove{
		XMLName:      "",
		ProfileToken: "MediaProfile000",
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{
				X:     -0.1, // (-1,1)
				Y:     0,    // (-1, 1)
				Space: "",
			},
			Zoom: onvif.Vector1D{
				X: 0, // +放大 -缩小, (0, 1)
			},
		},
		Timeout: "PT10S", // 移动持续时间
	})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(moveResponse)
	}

	// 设备不一定支持相对移动

	/*time.Sleep(time.Millisecond * 200)
	stopResponse, err := sdkptz.Call_Stop(ctx, dev, ptz.Stop{
		ProfileToken: "MediaProfile000",
		Zoom:         true,
	})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(stopResponse)
	}*/
}
