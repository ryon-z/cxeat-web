package router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"yelloment-api/config"
	"yelloment-api/controllers"
	"yelloment-api/database"
	"yelloment-api/global"
	"yelloment-api/middleware"
	"yelloment-api/utils"
	viewcontrollers "yelloment-api/view_controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter : setup router
func SetupRouter(dbType string) *gin.Engine {
	router := gin.Default()

	// We didn't use CORS yet
	/*
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{
			// envutil.GetGoDotEnvVariable("BASE_URL"),
		}
		corsConfig.AllowCredentials = true
		router.Use(cors.New(corsConfig))
	*/

	isTest := false
	if os.Getenv("IS_TEST") == "TRUE" {
		isTest = true
	}

	// Config 초기화
	config.InitStatusCodes()
	config.InitInjectionWords()

	// Utils 초기화
	utils.InitAllDowMap()

	// DB 세팅
	database.SetDB(dbType)

	// CssRandomVersion 초기화
	rand.Seed(time.Now().UnixNano())
	global.CssRandomVersion = strconv.Itoa(rand.Intn(10000000))

	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// 로그인
	router.POST("/login", authMiddleware.LoginHandler)
	router.NoRoute(func(c *gin.Context) {
		state := gin.H{"errorMessage": "없는 페이지 입니다."}
		global.UpdateGlobalState(c)
		utils.CombineCustomStateGlobalState(c, &state)

		c.HTML(http.StatusNotFound, "error.html", state)
	})

	auth := router.Group("/auth")
	auth.GET("/refresh", authMiddleware.RefreshHandler)

	// 유저
	user := router.Group("/user")
	user.GET("kakao", authMiddleware.LoginHandler)
	user.GET("naver", authMiddleware.LoginHandler)
	// if isTest {
	// 	user.POST("create", controllers.CreateUser)
	// }
	user.POST("create", controllers.CreateUser)
	user.Use(authMiddleware.MiddlewareFunc())
	{
		user.GET("info", controllers.ListOwnedUser)
		user.POST("edit", controllers.PartialUpdateOwnedUser)
	}

	// 주소
	address := router.Group("/address")
	address.GET("search", controllers.SearchAddress)
	address.Use(authMiddleware.MiddlewareFunc())
	{
		if isTest {
			address.POST("create", middleware.Generic("userAddress", "create", "AddressNo", []string{"RegDate"}))
		}
		address.GET("list", controllers.ListOwnedUserAddress)
		address.GET("retrieve/:id", controllers.RetrieveOwnedUserAddress) // client-side에서 ajax call 허용
		address.POST("add", controllers.CreateOwnedUserAddress)           // client-side에서 ajax call 허용
		address.POST("modify", controllers.PartialUpdateOwnedUserAddress) // client-side에서 ajax call 허용
		address.POST("remove", controllers.DeactivateOwnedUserAddress)    // client-side에서 ajax call 허용
	}

	// 카드
	card := router.Group("/card")
	card.Use(authMiddleware.MiddlewareFunc())
	{
		if isTest {
			card.POST("create", middleware.Generic("userCard", "create", "CardRegNo", []string{"RegDate"}))
		}
		card.POST("add", controllers.CreateOwnedUserCard)
		card.GET("list", controllers.ListOwnedUserCard)
		card.GET("retrieve/:id", controllers.RetrieveOwnedUserCard) // client-side에서 ajax call 허용
		card.POST("remove", controllers.DeactivateOwndUserCard)     // client-side에서 ajax call 허용
	}

	// 구독
	subs := router.Group("/subs")
	subs.Use(authMiddleware.MiddlewareFunc())
	{
		if isTest {
			subs.POST("create", middleware.Generic("subsMst", "create", "SubsNo", []string{"RegDate"}))
		}
		subs.GET("list", controllers.ListOwnedSubsMst)
		subs.POST("request", controllers.CreateOwnedSubsMst) // client-side에서 ajax call 허용
		subs.POST("pause", controllers.PauseOwnedSubsMst)
		subs.POST("cancel", controllers.DeactivateOwnedSubsMst)
		subs.POST("edit", controllers.PartialUpdateOwnedSubsMst)
	}

	// 주문
	order := router.Group("/order")
	order.Use(authMiddleware.MiddlewareFunc())
	{
		if isTest {
			order.POST("create", middleware.Generic("orderMst", "create", "OrderNo", []string{"RegDate"}))
		}
		order.POST("request/once", controllers.CreateOnceOwnedOrderMst)
		order.GET("list", controllers.ListOwnedOrder)
		order.POST("cancel", controllers.DeactivateOwnedOrder)
	}

	// 배너 조회
	banner := router.Group("/banner")
	if isTest {
		banner.POST("create", middleware.Generic("bannerMst", "create", "BannerNo", []string{"RegDate"}))
	}
	banner.GET("list", middleware.Generic("bannerMst", "list", "AddressNo", []string{"RegDate"}))

	// 약관 조회
	agreement := router.Group("/agreement")
	if isTest {
		agreement.POST("create", middleware.Generic("agreementMst", "create", "AgreementNo", []string{"RegDate"}))
	}
	agreement.GET("list", middleware.Generic("agreementMst", "list", "AddressNo", []string{"RegDate"}))

	// FAQ 조회
	faq := router.Group("/faq")
	if isTest {
		faq.POST("create", middleware.Generic("faqMst", "create", "FaqNo", []string{"RegDate"}))
	}
	faq.GET("list", middleware.Generic("faqMst", "list", "FaqNo", []string{"RegDate"}))

	// 태그 등록
	tag := router.Group("/tag")
	tag.Use(authMiddleware.MiddlewareFunc())
	{
		tag.POST("create", controllers.CreateOwnedTag)
		tag.POST("group/create", controllers.CreateOwnedTagGroup)
	}

	// 리뷰 등록
	review := router.Group("/review")
	review.POST("create", controllers.CreateReviewMst)
	review.POST("item/create", controllers.CreateItemReview)
	review.Use(authMiddleware.MiddlewareFunc())
	{
		review.POST("edit", controllers.PartialUpdateReviewMst)
		review.POST("item/edit", controllers.PartialUpdateItemReview)
	}

	// 구독 취소 알람
	alarm := router.Group("/alarm")
	alarm.Use(authMiddleware.MiddlewareFunc())
	{
		alarm.GET("subs/cancel", middleware.RenderController(controllers.AlarmCancelSubs))
	}

	// 주문 기록
	orderHist := router.Group("/order-hist")
	orderHist.Use(authMiddleware.MiddlewareFunc())
	{
		orderHist.POST("create", middleware.RenderController(controllers.CreateOrderHist))
	}

	// 구독 기록
	subsHist := router.Group("/subs-hist")
	subsHist.Use(authMiddleware.MiddlewareFunc())
	{
		subsHist.POST("create", middleware.RenderController(controllers.CreateSubsHist))
	}

	// Web page view
	// Load html template files
	templateDirPath := utils.GetWorkingDirPath() + "/views"
	var templatePaths []string
	getTemplatePaths(templateDirPath, &templatePaths)
	router.SetFuncMap(template.FuncMap{
		"modZero": func(dividen int, divisor int) bool {
			return dividen%divisor == 0
		},
		"minus": func(subtracted int, subtractor int) int {
			return subtracted - subtractor
		},
		"eqStrPtr": func(pointer *string, target string) bool {
			return *pointer == target
		},
	})
	router.LoadHTMLFiles(templatePaths...)
	router.StaticFS("/contents", http.Dir("contents"))
	router.StaticFile("/robots.txt", utils.GetWorkingDirPath()+"/robots.txt")
	router.StaticFile("/sitemap.xml", utils.GetWorkingDirPath()+"/sitemap.xml")
	router.StaticFile("/naver6bb63dc3f3e6a1076db40c00ae18849e.html", utils.GetWorkingDirPath()+"/naver6bb63dc3f3e6a1076db40c00ae18849e.html")
	router.GET("favicon.ico", func(c *gin.Context) {
		c.File("contents/images/favicon.ico")
	})

	router.GET("health-check", middleware.RenderWithParams(http.StatusOK, "health-check", gin.H{}))
	router.GET("", middleware.RenderController(viewcontrollers.HandleMain))
	web := router.Group("/web")
	web.GET("main", middleware.RenderController(viewcontrollers.HandleMain))

	web.GET("login", middleware.RenderWithParams(http.StatusOK, "login", gin.H{"navTitle": "로그인"}))
	web.GET("login/kakao", middleware.RenderController(viewcontrollers.HandleKakaoLogin))
	web.GET("login/naver", middleware.RenderController(viewcontrollers.HandleNaverLogin))

	web.GET("subs/intro", middleware.RenderWithParams(
		http.StatusOK, "subs-intro", gin.H{"navTitle": "요금제 안내", "activeNav": "subsIntro"}))
	web.GET("benefits", middleware.RenderWithParams(
		http.StatusOK, "benefits", gin.H{"navTitle": "혜택", "activeNav": "benefits"}))
	web.GET("story", middleware.RenderWithParams(
		http.StatusOK, "story", gin.H{"navTitle": "브랜드 스토리", "activeNav": "story"}))

	web.GET("review/insert/once", middleware.RenderController(viewcontrollers.HandleInsertReviewOnce))
	web.Use(authMiddleware.MiddlewareFunc())
	{
		web.GET("survey/completed", middleware.RenderController(viewcontrollers.HandleSurveyCompleted))
		web.GET("survey/result", middleware.RenderController(viewcontrollers.HandleSurveyResult))
		web.GET("survey/main", middleware.RenderController(viewcontrollers.HandleSurveyMain))
		web.POST("survey/order-and-delivery/info", middleware.RenderController(viewcontrollers.HandleOrderAndDeliveryInfo))
		web.GET("delivery/manage", middleware.RenderController(viewcontrollers.HandleManageDelivery))
		web.GET("delivery/insert", middleware.RenderController(viewcontrollers.HandleInsertDelivery))
		web.GET("delivery/modify", middleware.RenderController(viewcontrollers.HandleModifyDelivery))
		web.GET("card/manage", middleware.RenderController(viewcontrollers.HandleManageCard))
		web.GET("card/insert", middleware.RenderController(viewcontrollers.HandleInsertCard))
		web.GET("my-page/main", middleware.RenderController(viewcontrollers.HandleMyPageMain))
		web.GET("my-page/user/info", middleware.RenderController(viewcontrollers.HandleMyPageUserInfo))
		web.GET("my-page/subs/list", middleware.RenderController(viewcontrollers.HandleMyPageSubsList))
		web.GET("my-page/subs/info", middleware.RenderController(viewcontrollers.HandleMyPageSubsInfo))
		web.GET("my-page/subs/edit", middleware.RenderController(viewcontrollers.HandleEditSubs))
		web.GET("my-page/subs/cancel", middleware.RenderWithParams(http.StatusOK, "my-page/cancel-subs", gin.H{"navTitle": "구독해지"}))
		web.GET("my-page/order/history", middleware.RenderController(viewcontrollers.HandleMyPageOrderHistory))
		web.GET("my-page/order/detail", middleware.RenderController(viewcontrollers.HandleMyPageOrderDetail))
		web.GET("my-page/review/list", middleware.RenderController(viewcontrollers.HandleMyPageReviewList))
		web.GET("my-page/review/detail", middleware.RenderController(viewcontrollers.HandleMyPageReviewDetail))
		web.GET("my-page/review/detail/edit", middleware.RenderController(viewcontrollers.HandleModifyReview))
		web.GET("my-page/withdrawal", middleware.RenderController(viewcontrollers.HandleMyPageWithdrawal))

		web.GET("signup", middleware.RenderController(viewcontrollers.HandleSignup))
		web.GET("refresh", controllers.Refresh)
		web.GET("logout", viewcontrollers.Logout)
	}

	return router
}

// getTemplatePaths : template 파일 경로를 resultPaths에 담음
func getTemplatePaths(dirPath string, resultPaths *[]string) bool {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", dirPath, file.Name())
		if file.IsDir() {
			getTemplatePaths(path, resultPaths)
		} else {
			*resultPaths = append(*resultPaths, path)
		}
	}

	return true
}
