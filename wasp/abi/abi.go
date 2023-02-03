package abi

import (
	"github.com/berachain/stargazer/wasp/models"
	"gorm.io/gorm"
)

func CreateDefaultAbi(db *gorm.DB) {
	erc20 := &models.Abi{
		Abi: ERC20ABI,
		Tag: "ERC20",
	}

	erc721 := &models.Abi{
		Abi: ERC721ABI,
		Tag: "ERC721",
	}

	db.Create(&erc20)
	db.Create(&erc721)
}
