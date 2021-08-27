package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e421083458/filecoin_sectors/dao"
	"github.com/e421083458/filecoin_sectors/dto"
	"github.com/e421083458/filecoin_sectors/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

func LoopSector() bool {
	//c := &gin.Context{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		return false
	}
	//添加请求头
	headers := http.Header{"Authorization": []string{"Bearer " + public.AuthToken}}
	var api lotusapi.StorageMinerStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+public.Addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Printf("connecting with lotus failed: %s", err)
		return false
	}
	defer closer()
	//获取扇区列表
	sectorsList, sectorsListErr := api.SectorsList(context.Background())
	if sectorsListErr != nil {
		log.Printf("calling SectorsList: %s", sectorsListErr)
	}
	//获取所有扇区文件位置
	storageLocal, storageLocalErr := api.StorageLocal(context.Background())
	if storageLocalErr != nil {
		log.Printf("calling StorageLocal: %s", err)
	}
	for _, val := range sectorsList {
		//获取扇区详情
		sectorsStatus, sectorsStatusErr := api.SectorsStatus(context.Background(), val, true)
		if sectorsStatusErr != nil {
			log.Printf("calling SectorsStatus: %s", sectorsStatusErr)
		}
		fmt.Println(sectorsStatus.State)
		fmt.Println(sectorsStatus.Expiration)
		fmt.Println(sectorsStatus.PreCommitMsg)
		fmt.Println(sectorsStatus.CommitMsg)
		//扇区状态非proving则跳过 、无commit消息跳过、无PreCommi跳过
		if !strings.EqualFold(string(sectorsStatus.State), "Proving") {
			continue
		}
		if sectorsStatus.PreCommitMsg == nil {
			continue
		}
		if sectorsStatus.CommitMsg == nil {
			continue
		}
		sid := abi.SectorID{
			Miner:  abi.ActorID(public.Mid),
			Number: abi.SectorNumber(val),
		}
		//查找扇区 获取扇区id 、Precommit 、Commit
		storageFindSector, storageFindSectorErr := api.StorageFindSector(context.Background(), sid, 1, 34359738368, true)
		if storageFindSectorErr != nil {
			log.Printf("calling SectorsStatus: %s", sectorsStatusErr)
		}
		id := storageFindSector[0].ID
		fileurl := storageFindSector[0].URLs
		locationUrl := storageLocal[id]
		/*
			PreCommitSector  ==   Precommit
			ProveCommitSector ==  Commit
		*/
		//获取PreCommitSector质押数量
		preCommitMsg := sectorsStatus.PreCommitMsg
		sprintfpreCommitMsg := fmt.Sprintf("%d", preCommitMsg)
		preCommitRs := get(public.FirfoxUrl + sprintfpreCommitMsg)
		preCommitInfo := dto.FirFoxInfo{}
		json.Unmarshal([]byte(preCommitRs), &preCommitInfo)
		//获取ProveCommitSector质押币数量
		CommitMsg := sectorsStatus.CommitMsg
		sprintfCommitMsg := fmt.Sprintf("%d", CommitMsg)
		commitRs := get(public.FirfoxUrl + sprintfCommitMsg)
		commitInfo := dto.FirFoxInfo{}
		json.Unmarshal([]byte(commitRs), &commitInfo)
		preCommitFil := preCommitInfo.Value.Div(decimal.NewFromFloat(math.Pow10(18)))
		proveCommitFil := commitInfo.Value.Div(decimal.NewFromFloat(math.Pow10(18)))
		totalPledgeFil := preCommitFil.Add(proveCommitFil)
		sector := &dao.Sectors{
			SectorId:       int(val),
			Nonce:          commitInfo.Nonce,
			SectorStatus:   string(sectorsStatus.State),
			Expiration:     int(sectorsStatus.Expiration),
			ExpirationStr:  "",
			LocationUrl:    locationUrl,
			FileUrl:        fileurl[0],
			PreCommitFil:   preCommitFil,
			ProveCommitFil: proveCommitFil,
			TotalPledgeFil: totalPledgeFil,
		}
		sectorInfo, sectorsInfoErr := sector.FindBySectorId(nil, tx, int64(val))
		if sectorsInfoErr != nil {
			log.Printf("dao FindBySectorId Err:%s", sectorsInfoErr)
		}
		if sectorInfo != nil {
			//update
			//sector.Update()
			log.Println("update")
			sector.Update(nil, tx, int64(val))
		} else {
			//save
			log.Println("save")
			sector.Save(nil, tx)
		}
	}
	return true
}

/**
发送get请求，返回str结果
*/
func get(url string) (rs string) {
	rs = ""
	resp, err := http.Get(url)
	if err != nil {
		log.Println("htpp get url: %s, err : %s", url, err)
		return ""
	}
	defer resp.Body.Close()
	respString, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		log.Println("ioutil.ReadAll url: %s, err: %s ", url, respErr)
	}
	return string(respString)
}
