package handlers

import (
	"PRACTICA_Spulling_Lpulling/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"time"
)

func ObtenerUsuarios(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		filas, err := db.Query("SELECT ID, nombre FROM usuarios")
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al consultar la base de datos"})
			return
		}
		defer filas.Close()

		var usuarios []models.Usuario
		for filas.Next() {
			var u models.Usuario
			if err := filas.Scan(&u.ID, &u.Nombre); err != nil {
				c.JSON(500, gin.H{"error": "Error al leer los datos"})
				return
			}
			usuarios = append(usuarios, u)
		}

		c.JSON(200, usuarios)
	}
}

func CrearUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u models.Usuario
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Datos inválidos"})
			return
		}

		_, err := db.Exec("INSERT INTO usuarios (nombre) VALUES (?)", u.Nombre)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al crear el usuario"})
			return
		}

		c.JSON(201, gin.H{"mensaje": "Usuario creado"})
	}
}

func ActualizarUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u models.Usuario
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Datos inválidos"})
			return
		}

		_, err := db.Exec("UPDATE usuarios SET nombre = ? WHERE ID = ?", u.Nombre, u.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al actualizar el usuario"})
			return
		}

		c.JSON(200, gin.H{"mensaje": "Usuario actualizado"})
	}
}

func EliminarUsuario(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u models.Usuario
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Datos inválidos"})
			return
		}

		_, err := db.Exec("DELETE FROM usuarios WHERE ID = ?", u.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar el usuario"})
			return
		}

		c.JSON(200, gin.H{"mensaje": "Usuario eliminado"})
	}
}

func SincronizarReplica(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		filas, err := db.Query("SELECT ID, nombre FROM usuarios")
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al consultar la base de datos"})
			return
		}
		defer filas.Close()

		var usuarios []models.Usuario
		for filas.Next() {
			var u models.Usuario
			if err := filas.Scan(&u.ID, &u.Nombre); err != nil {
				c.JSON(500, gin.H{"error": "Error al leer los datos"})
				return
			}
			usuarios = append(usuarios, u)
		}

		c.JSON(200, usuarios)
	}
}

func ShortPolling(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ultimoID int
		err := db.QueryRow("SELECT COALESCE(MAX(ID), 0) FROM usuarios").Scan(&ultimoID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al obtener el último ID"})
			return
		}

		time.Sleep(2 * time.Second)

		var nuevosUsuarios []models.Usuario
		filas, err := db.Query("SELECT ID, nombre FROM usuarios WHERE ID > ?", ultimoID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al consultar la base de datos"})
			return
		}
		defer filas.Close()

		for filas.Next() {
			var u models.Usuario
			if err := filas.Scan(&u.ID, &u.Nombre); err != nil {
				c.JSON(500, gin.H{"error": "Error al leer los datos"})
				return
			}
			nuevosUsuarios = append(nuevosUsuarios, u)
		}

		if len(nuevosUsuarios) > 0 {
			c.JSON(200, gin.H{"mensaje": "Nuevos usuarios encontrados", "usuarios": nuevosUsuarios})
		} else {
			c.JSON(200, gin.H{"mensaje": "No hay nuevos usuarios"})
		}
	}
}

func LongPolling(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ultimoID int
		err := db.QueryRow("SELECT COALESCE(MAX(ID), 0) FROM usuarios").Scan(&ultimoID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error al obtener el último ID"})
			return
		}

		timeout := time.After(30 * time.Second)
		tick := time.Tick(2 * time.Second)

		for {
			select {
			case <-timeout:
				c.JSON(200, gin.H{"mensaje": "No hay cambios"})
				return
			case <-tick:
				var nuevosUsuarios []models.Usuario
				filas, err := db.Query("SELECT ID, nombre FROM usuarios WHERE ID > ?", ultimoID)
				if err != nil {
					c.JSON(500, gin.H{"error": "Error al consultar la base de datos"})
					return
				}
				defer filas.Close()

				for filas.Next() {
					var u models.Usuario
					if err := filas.Scan(&u.ID, &u.Nombre); err != nil {
						c.JSON(500, gin.H{"error": "Error al leer los datos"})
						return
					}
					nuevosUsuarios = append(nuevosUsuarios, u)
				}

				if len(nuevosUsuarios) > 0 {
					c.JSON(200, gin.H{"mensaje": "Nuevos usuarios encontrados", "usuarios": nuevosUsuarios})
					return
				}
			}
		}
	}
}
