func CORSMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

if c.Request.Method == "OPTIONS" {
c.AbortWithStatus(204)
return
}

c.Next()
}
}

func main() {
r := gin.Default()
r.Use(CORSMiddleware()) // Gunakan middleware sebelum mendefinisikan route

// ... route kamu
r.Run(":8080")
}