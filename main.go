package main

import (
	"PRACTICA_Spulling_Lpulling/handlers"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {

	dbPrincipal, err := sql.Open("mysql", "root:roooooot@tcp(127.0.0.1:3306)/principal_db")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos principal:", err)
	}
	defer dbPrincipal.Close()

	dbReplica, err := sql.Open("mysql", "root:roooooot@tcp(127.0.0.1:3306)/replica_db")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos de réplica:", err)
	}
	defer dbReplica.Close()

	r := gin.Default()

	r.GET("/usuarios", handlers.ObtenerUsuarios(dbPrincipal))
	r.POST("/usuarios/crear", handlers.CrearUsuario(dbPrincipal))
	r.PUT("/usuarios/actualizar", handlers.ActualizarUsuario(dbPrincipal))
	r.DELETE("/usuarios/eliminar", handlers.EliminarUsuario(dbPrincipal))
	r.GET("/replica/sincronizar", handlers.SincronizarReplica(dbPrincipal))

	r.GET("/replica/usuarios", handlers.ObtenerUsuarios(dbReplica))

	r.GET("/short-polling", handlers.ShortPolling(dbPrincipal))
	r.GET("/long-polling", handlers.LongPolling(dbPrincipal))

	log.Println("Servidor principal y réplica iniciados en :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
