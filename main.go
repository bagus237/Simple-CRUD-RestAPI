package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)
//Bismillah

var db *gorm.DB
//Koneksi ke MySQL
func init() {
	var err error
	db, err =
		gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/service?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Gagal Conect Ke Database")
	}
	db.AutoMigrate(&Person{})
}

//Pembuatan record dan field aduh salah
type (
	Person struct {
		Name        string `json:"name"`
		Addres      string `json:"addres"`
		Phone_Number string `json:"phone_number"`
		Gender string `json:"gender"`
		gorm.Model
	}
	//Untuk melakukan update field
	editPerson struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Addres     string `json:"addres"`
		Gender      string   `json:"gender"`
		Phone_Number string `json:"phone_number"`
		
	}
)
// Untuk Creat Tabel
func cretedPerson(c *gin.Context) {
	var std editPerson
	var model Person
	c.Bind(&std)
	validasi := validatorCreated(std)
	model = transferEPToModel(std)
	if validasi != "" {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": validasi})
	} else {
		db.Create(&model)
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": model})
	}
}
//Untuk Select Tabel
func fetchAllPerson(c *gin.Context) {
	var model [] Person
	var EP [] editPerson

	db.Find(&model)

	if len(model) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}

	for _, item := range model {
		EP = append(EP, transferModelToEP(item))
	}
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": EP})
}
//Untuk Select Single Person
func fetchSinglePerson(c *gin.Context) {
	var model Person
	var EP editPerson

	modelID := c.Param("id")
	db.Find(&model, modelID)

	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}
	EP = transferModelToEP(model)
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": EP})
}
// Untuk Update Person
func updatePerson(c *gin.Context) {
	var model Person
	var EP editPerson
	modelID := c.Param("id")
	db.First(&model, modelID)

	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Data Tidak Ada"})
	}
	c.Bind(&EP)

	validasi := validatorCreated(EP)
	if validasi != "" {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": validasi})
	} else {
		db.Model(&model).Update(transferEPToModel(EP))
		c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": model})
	}
}
//Untuk Delete Person
func deletePerson(c *gin.Context) {
	var model Person
	modelID := c.Param("id")

	db.First(&model, modelID)
	if model.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusNotFound, "result": "Datat Tidak di Temukan"})
	}
	db.Delete(model)
	c.JSON(http.StatusOK, gin.H{"message": http.StatusOK, "result": "Data Telah berhasil di hapus"})
}
//Untuk Update Persone
func transferModelToEP(model Person) editPerson {
	var EP editPerson
	EP = editPerson{
		ID:          model.ID,
		Name:        model.Name,
		Gender: 	 model.Gender,
		Addres:      model.Addres,
		Phone_Number: model.Phone_Number,
	
	
	}
	return EP
}

func transferEPToModel(EP editPerson) Person {
	var model Person
	model = Person{
		Name:        EP.Name,
		Gender:      EP.Gender,
		Addres:      EP.Addres,
		Phone_Number: EP.Phone_Number,
	}
	return model
}
// Untuk Validasi data not null
func validatorCreated(EP editPerson) string {

	var kosong string = " Is Empty"

	if EP.Name == "" {
		return "name" + kosong
	}

	if EP.Addres == "" {
		return "addres" + kosong
	}

	if EP.Gender == "" {
		return "gender" + kosong
	}

	if EP.Phone_Number == "" {
		return "phone_number" + kosong
	}

	return ""
}

func main() {
//Untuk Alamat Urlnya
	router := gin.Default()
	v1 := router.Group("/api/service")
	{
		v1.POST("", cretedPerson)
		v1.GET("", fetchAllPerson)
		v1.GET("/:id", fetchSinglePerson)
		v1.PUT("/:id", updatePerson)
		v1.DELETE("/:id", deletePerson)
	}
	router.Run(":2000")
}
