package handler

import (
	"elichika/config"
	"elichika/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func UpdateCardNewFlag(ctx *gin.Context) {
	signBody, _ := sjson.Set(GetData("updateCardNewFlag.json"),
		"user_model_diff.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ChangeIsAwakeningImage(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.CardAwakeningReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	loginData := GetUserData("userCard.json")
	cardInfo := model.CardInfo{}
	gjson.Parse(loginData).Get("user_card_by_card_id").
		ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &cardInfo); err != nil {
					panic(err)
				}

				if cardInfo.CardMasterID == req.CardMasterID {
					cardInfo.IsAwakeningImage = req.IsAwakeningImage

					k := "user_card_by_card_id." + key.String() + ".is_awakening_image"
					SetUserData("userCard.json", k, req.IsAwakeningImage)

					return false
				}
			}
			return true
		})

	userCardInfo := []any{}
	userCardInfo = append(userCardInfo, cardInfo.CardMasterID)
	userCardInfo = append(userCardInfo, cardInfo)

	// Update user profile
	cardMasterId := gjson.Parse(GetUserData("fetchProfile.json")).Get("profile_info.basic_info.recommend_card_master_id").Int()
	if cardMasterId == int64(req.CardMasterID) {
		SetUserData("fetchProfile.json", "profile_info.basic_info.is_recommend_card_image_awaken", req.IsAwakeningImage)
	}

	cardResp := GetData("changeIsAwakeningImage.json")
	cardResp, _ = sjson.Set(cardResp, "user_model_diff.user_status", GetUserStatus())
	cardResp, _ = sjson.Set(cardResp, "user_model_diff.user_card_by_card_id", userCardInfo)
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ChangeFavorite(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.CardFavoriteReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	cardData := GetUserData("userCard.json")
	cardInfo := model.CardInfo{}
	gjson.Parse(cardData).Get("user_card_by_card_id").
		ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &cardInfo); err != nil {
					panic(err)
				}

				if cardInfo.CardMasterID == req.CardMasterID {
					cardInfo.IsFavorite = req.IsFavorite

					k := "user_card_by_card_id." + key.String() + ".is_favorite"
					SetUserData("userCard.json", k, req.IsFavorite)

					return false
				}
			}
			return true
		})

	userCardInfo := []any{}
	userCardInfo = append(userCardInfo, cardInfo.CardMasterID)
	userCardInfo = append(userCardInfo, cardInfo)

	cardResp := GetData("changeFavorite.json")
	cardResp, _ = sjson.Set(cardResp, "user_model_diff.user_status", GetUserStatus())
	cardResp, _ = sjson.Set(cardResp, "user_model_diff.user_card_by_card_id", userCardInfo)
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func GetOtherUserCard(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	userCardReq := model.UserCardReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &userCardReq); err != nil {
		panic(err)
	}
	// fmt.Println(liveStartReq)

	var newUserCardInfo model.NewCardInfo
	var cardInfo string
	partnerList := gjson.Parse(GetData("fetchLivePartners.json")).Get("partner_select_state.live_partners")
	partnerList.ForEach(func(k, v gjson.Result) bool {
		userId := v.Get("user_id").Int()
		if userId == userCardReq.UserID {
			v.Get("card_by_category").ForEach(func(kk, vv gjson.Result) bool {
				if vv.IsObject() {
					cardId := vv.Get("card_master_id").Int()
					if cardId == userCardReq.CardMasterID {
						cardInfo = vv.String()
						// fmt.Println(cardInfo)
						return false
					}
				}
				return true
			})
			return false
		}
		return true
	})

	if err := json.Unmarshal([]byte(cardInfo), &newUserCardInfo); err != nil {
		panic(err)
	}

	userCardResp := GetData("getOtherUserCard.json")
	userCardResp, _ = sjson.Set(userCardResp, "other_user_card", newUserCardInfo)
	resp := SignResp(ctx.GetString("ep"), userCardResp, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchTrainingTree(ctx *gin.Context) {
	signBody := GetData("fetchTrainingTree.json")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func GradeUpCard(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	fmt.Println(reqBody.String())

	var req model.GradeUpCardReq
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}
	fmt.Println(req)

	var cardInfo model.CardInfo
	gjson.Parse(GetUserData("userCard.json")).Get("user_card_by_card_id").
		ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() && value.Get("card_master_id").Int() == req.CardMasterID {
				keyCard := "user_card_by_card_id." + key.String()
				grade := gjson.Get(value.String(), "grade").Int()
				keyCardInfo, _ := sjson.Set(value.String(), "grade", grade+1)
				if err := json.Unmarshal([]byte(keyCardInfo), &cardInfo); err != nil {
					panic(err)
				}
				//
				SetUserData("userCard.json", keyCard, cardInfo)

				return false
			}
			return true
		})

	var userCard []any
	userCard = append(userCard, req.CardMasterID)
	userCard = append(userCard, cardInfo)

	var memberInfo model.UserMemberInfo
	memberId := GetMemberMasterIdByCardMasterId(int(req.CardMasterID))
	gjson.Parse(GetUserData("memberSettings.json")).Get("user_member_by_member_id").
		ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() && int(value.Get("member_master_id").Int()) == memberId {
				keyMemberInfo, _ := sjson.Set(value.String(), "love_point_limit", 13181880)
				keyMemberInfo, _ = sjson.Set(keyMemberInfo, "love_level", 500)
				if err := json.Unmarshal([]byte(keyMemberInfo), &memberInfo); err != nil {
					panic(err)
				}
				return false
			}
			return true
		})

	var memberData []any
	memberData = append(memberData, memberId)
	memberData = append(memberData, memberInfo)

	triggerId := time.Now().UnixNano()
	triggerInfo := model.CardGradeUpTriggerInfo{
		TriggerID:            triggerId,
		CardMasterID:         req.CardMasterID,
		BeforeLoveLevelLimit: 497,
		AfterLoveLevelLimit:  500,
	}
	var gradeUpTrigger []any
	gradeUpTrigger = append(gradeUpTrigger, triggerId)
	gradeUpTrigger = append(gradeUpTrigger, triggerInfo)

	signBody := GetData("gradeUpCard.json")
	signBody, _ = sjson.Set(signBody, "user_model_diff.user_member_by_member_id", memberData)
	signBody, _ = sjson.Set(signBody, "user_model_diff.user_card_by_card_id", userCard)
	signBody, _ = sjson.Set(signBody, "user_model_diff.user_info_trigger_card_grade_up_by_trigger_id", gradeUpTrigger)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ReadCardGradeUp(ctx *gin.Context) {
	triggerId := gjson.Parse(ctx.GetString("reqBody")).Array()[0].Get("trigger_id").Int()
	fmt.Println(triggerId)

	var triggerInfo []any
	triggerInfo = append(triggerInfo, triggerId)
	triggerInfo = append(triggerInfo, nil)

	signBody := GetData("readCardGradeUp.json")
	signBody, _ = sjson.Set(signBody, "user_model.user_info_trigger_card_grade_up_by_trigger_id", triggerInfo)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
