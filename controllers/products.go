package controllers

import (
	// "html/template"

	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Addproducts(c *gin.Context) {
	c.HTML(http.StatusOK, "addProduct.html", nil)
}
func PostAddproducts(c *gin.Context) {
	var err error

	name := c.Request.FormValue("productName")
	description := c.Request.FormValue("productDescription")
	stock, _ := strconv.Atoi(c.Request.FormValue("productStock"))
	price, _ := strconv.Atoi(c.Request.FormValue("productPrice"))
	category_name := c.Request.FormValue("categoryName")
	brand_name := c.Request.FormValue("brandName")

	//Get the image file
	file, err := c.FormFile("productImage")
	if err != nil {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to upload image",
		})
		return
	}
	//Save the image file
	err = c.SaveUploadedFile(file, "./static/images/"+file.Filename)
	if err != nil {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to save image",
		})
		return
	}
	if err != nil {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to add product",
		})
		return
	}
	var dtproduct models.Product
	database.DB.Where("name=?", name).First(&dtproduct)

	if name == dtproduct.Name {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "This product already exist",
		})
		return
	}
	//----------------

	//adding category
	var dtcategory models.Category
	database.DB.Table("categories").Where("name=?", category_name).Scan(&dtcategory)
	if dtcategory.Name != category_name {
		log.Println("category is not exist")
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "This category is not exist",
		})
		return
	}
	//----------------

	result := database.DB.Create(&models.Product{
		Name:          name,
		Description:   description,
		Stock:         stock,
		Price:         price,
		Category_Name: category_name,
		Brand_Name:    brand_name,
		Image:         file.Filename,
	})
	if result.Error != nil {
		log.Println("Failed to add product", err)
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to add product",
		})
		return
	}
	c.HTML(http.StatusOK, "addProduct.html", gin.H{
		"message": "Successfully add product",
	})
	c.Redirect(http.StatusSeeOther, "/user/products-list")
}
func AdminListproducts(c *gin.Context) {
	type products struct {
		Id            uint
		Name          string `gorm:"not null"`
		Description   string `gorm:"not null"`
		Stock         int    `gorm:"not null"`
		Price         int    `gorm:"not null"`
		Category_Name string `gorm:"not null"`
		Brand_Name    string `gorm:"not null"`
		Image         string `gorm:"not null"`
	}
	var product []products
	result := database.DB.Table("products").Select("id,name,description,stock,price,category_name,brand_name,image").Scan(&product)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "productsList2.html", gin.H{
			"error": "Failed to list product",
		})
		return
	}
	c.HTML(200, "adminProductlist.html", product)
}
func Listproducts(c *gin.Context) {
	data := DtTables()
	c.HTML(200, "productsList2.html", data)
}

func ProductDetails(c *gin.Context) {
	productiD := c.Param("id")

	sizepd, _ := strconv.Atoi(c.PostForm("size"))
	fmt.Println("sizepd : ", sizepd)

	var products models.Product
	database.DB.Table("products").Where("id=?", productiD).First(&products)
	c.HTML(http.StatusOK, "productDetails2.html", products)
}

func Editproduct(c *gin.Context) {
	c.HTML(200, "editproduct.html", nil)
}
func PostEditproduct(c *gin.Context) {
	var err error

	name := c.Request.FormValue("productName")
	description := c.Request.FormValue("productDescription")
	stock, _ := strconv.Atoi(c.Request.FormValue("productStock"))
	price, _ := strconv.Atoi(c.Request.FormValue("productPrice"))
	category_name := c.Request.FormValue("categoryName")
	brand_name := c.Request.FormValue("brandName")

	//Get the image file
	file, err := c.FormFile("productImage")
	if err != nil {
		c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
			"error": "Failed to upload image",
		})
		return
	}
	//Save the image file
	err = c.SaveUploadedFile(file, "./static/images/"+file.Filename)
	if err != nil {
		c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
			"error": "Failed to save the edited image",
		})
		return
	}
	//checking category
	var dtcategory models.Category
	database.DB.Table("categories").Where("name=?", category_name).Scan(&dtcategory)
	if dtcategory.Name != category_name {
		log.Println("category is not exist")
		c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
			"error": "This category is not exist",
		})
		return
	}

	iD, _ := strconv.Atoi(c.Param("id"))

	result := database.DB.Model(&models.Product{}).Where("id=?", iD).Updates(map[string]interface{}{
		"name":          name,
		"description":   description,
		"stock":         stock,
		"price":         price,
		"category_name": category_name,
		"brand_name":    brand_name,
		"image":         file.Filename,
	})
	if result.Error != nil {
		log.Println("Failed to edit product")
		c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
			"error": "Failed to edit product",
		})
		return
	}
	var product models.Product
	database.DB.Table("products").Where("id=?", iD).First(&product)
	c.HTML(http.StatusOK, "editproduct.html", product)
	c.Redirect(http.StatusSeeOther, "/admin-products-list")
}
func Deleteproduct(c *gin.Context) {

	type products struct {
		Id            uint
		Name          string `gorm:"not null"`
		Description   string `gorm:"not null"`
		Stock         int    `gorm:"not null"`
		Price         int    `gorm:"not null"`
		Category_Name string `gorm:"not null"`
		Brand_Name    string `gorm:"not null"`
		Image         string `gorm:"not null"`
	}
	var product []products
	result := database.DB.Table("products").Select("id,name,description,stock,price,category_name,brand_name,image").Scan(&product)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "productsList2.html", gin.H{
			"error": "Failed to list product",
		})
		return
	}
	c.HTML(200, "adminProductlist.html", product)

	idstr := c.Param("prdctid")
	id, _ := strconv.Atoi(idstr)

	res := database.DB.Where("id", id).Delete(&models.Product{})
	if res.RowsAffected == 0 {
		c.HTML(400, "adminProductlist.html", gin.H{
			"error": "Failed to find product",
		})
		return
	}
	// c.HTML(http.StatusOK, "adminProductlist.html", product)
	c.Redirect(303, "/admin-products-list")
}
func DtTables() interface{} {

	var products []models.Product
	var categories []models.Category
	var carts []models.Cart
	var addresses []models.Address
	var users []models.User
	var brands []models.Brand

	database.DB.Find(&products)
	database.DB.Find(&categories)
	database.DB.Find(&carts)
	database.DB.Find(&addresses)
	database.DB.Find(&users)
	database.DB.Find(&brands)

	data := struct {
		Products   []models.Product
		Categories []models.Category
		Carts      []models.Cart
		Addresses  []models.Address
		Users      []models.User
		Brands     []models.Brand
	}{
		Products:   products,
		Categories: categories,
		Carts:      carts,
		Addresses:  addresses,
		Users:      users,
		Brands:     brands,
	}

	return data
}