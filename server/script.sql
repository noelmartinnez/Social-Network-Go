-- Crear la tabla USUARIO
CREATE TABLE USUARIO (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    hashed_username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    nombre VARCHAR(100),
    apellidos VARCHAR(100)
);

-- Crear la tabla POST
CREATE TABLE POST (
    id INT AUTO_INCREMENT PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    texto VARCHAR(255) NOT NULL,
    usuario_id INT,
    FOREIGN KEY (usuario_id) REFERENCES USUARIO(id)
);

-- Crear la tabla MENSAJE
CREATE TABLE MENSAJE (
    id INT AUTO_INCREMENT PRIMARY KEY,
    remitente_id INT NOT NULL,
    destinatario_id INT NOT NULL,
    contenido TEXT NOT NULL,
    fecha_envio TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (remitente_id) REFERENCES USUARIO(id),
    FOREIGN KEY (destinatario_id) REFERENCES USUARIO(id)
);

-- Crear la tabla de ROLES
CREATE TABLE ROLES (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

-- Insertar roles predeterminados
INSERT INTO ROLES (nombre) VALUES ('usuario');
INSERT INTO ROLES (nombre) VALUES ('moderador');
INSERT INTO ROLES (nombre) VALUES ('administrador');

-- Modificar la tabla USUARIO para incluir el ID del rol
ALTER TABLE USUARIO ADD COLUMN rol_id INT;
ALTER TABLE USUARIO ADD FOREIGN KEY (rol_id) REFERENCES ROLES(id);

CREATE TABLE `event_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `event_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` int DEFAULT NULL,
  `event_type` varchar(50) DEFAULT NULL,
  `event_description` text,
  `ip_address` varchar(45) DEFAULT NULL,
  `additional_data` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `event_time` (`event_time`)
);

-- Crear la tabla GRUPO
CREATE TABLE GRUPO (
    id INT AUTO_INCREMENT PRIMARY KEY,
    descripcion VARCHAR(255) NOT NULL,
    administrador_id INT NOT NULL,
    FOREIGN KEY (administrador_id) REFERENCES USUARIO(id)
);

-- Añadir la columna grupo_id a la tabla USUARIO
ALTER TABLE USUARIO ADD COLUMN grupo_id INT;
ALTER TABLE USUARIO ADD FOREIGN KEY (grupo_id) REFERENCES GRUPO(id);

-- Cambiar el delimitador
DELIMITER //

-- Añadir la restricción de que un grupo puede tener máximo 5 usuarios
CREATE TRIGGER grupo_max_usuarios BEFORE INSERT ON USUARIO
FOR EACH ROW
BEGIN
    IF (NEW.grupo_id IS NOT NULL) THEN
        IF ((SELECT COUNT(*) FROM USUARIO WHERE grupo_id = NEW.grupo_id) >= 5) THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Un grupo no puede tener más de 5 usuarios';
        END IF;
    END IF;
END;
//

CREATE TRIGGER grupo_max_usuarios_update BEFORE UPDATE ON USUARIO
FOR EACH ROW
BEGIN
    IF (NEW.grupo_id IS NOT NULL) THEN
        IF ((SELECT COUNT(*) FROM USUARIO WHERE grupo_id = NEW.grupo_id AND id != NEW.id) >= 5) THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Un grupo no puede tener más de 5 usuarios';
        END IF;
    END IF;
END;
//

-- Cambiar el delimitador de vuelta a punto y coma
DELIMITER ;

-- Modificar la tabla POST para incluir una columna grupo_id y una columna es_publico
ALTER TABLE POST ADD COLUMN grupo_id INT;
ALTER TABLE POST ADD COLUMN es_publico BOOLEAN DEFAULT TRUE;
ALTER TABLE POST ADD FOREIGN KEY (grupo_id) REFERENCES GRUPO(id);
