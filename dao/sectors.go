package dao

import (
	"github.com/e421083458/filecoin_sectors/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"time"
)

type Sectors struct {
	Id             int64           `json:"id" gorm:"primary_key"`
	SectorId       int             `json:"sectorId" gorm:"column:sector_id"`
	Nonce          int             `json:"nonce" gorm:"column:nonce"`
	SectorStatus   string          `json:"sectorStatus" gorm:"column:sector_status"`
	Expiration     int             `json:"expiration" gorm:"column:expiration"`
	ExpirationStr  string          `json:"expirationStr" gorm:"column:expiration_str"`
	FileUrl        string          `json:"fileUrl" gorm:"column:file_url"`
	LocationUrl    string          `json:"locationUrl" gorm:"column:location_url"`
	PreCommitFil   decimal.Decimal `json:"preCommitFil" gorm:"column:pre_commit_fil"`
	ProveCommitFil decimal.Decimal `json:"proveCommitFil" gorm:"column:prove_commit_fil"`
	TotalPledgeFil decimal.Decimal `json:"totalPledgeFil" gorm:"column:total_pledge_fil"`
	CreatedAt      time.Time       `json:"createTime" gorm:"column:create_time"`
	UpdatedAt      time.Time       `json:"updateTime" gorm:"column:update_time"`
}

func (f *Sectors) TableName() string {
	return "t_sectors"
}

func (f *Sectors) FindBySectorId(c *gin.Context, tx *gorm.DB, sectorId int64) (*Sectors, error) {
	model := &Sectors{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where("sector_id = ?", sectorId).Find(model).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (f *Sectors) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Create(f).Error; err != nil {
		return err
	}
	return nil
}

func (f *Sectors) Update(c *gin.Context, tx *gorm.DB, sectorid int64) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Model(f).Where("sector_id = ?", sectorid).Updates(f).Error
	if err != nil {
		return err
	}
	return nil
}
