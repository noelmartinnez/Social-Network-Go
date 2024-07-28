-- MySQL dump 10.13  Distrib 8.0.31, for Win64 (x86_64)
--
-- Host: localhost    Database: sds
-- ------------------------------------------------------
-- Server version	8.0.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `event_logs`
--

DROP TABLE IF EXISTS `event_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
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
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `event_logs`
--

LOCK TABLES `event_logs` WRITE;
/*!40000 ALTER TABLE `event_logs` DISABLE KEYS */;
INSERT INTO `event_logs` VALUES (1,'2024-05-25 15:19:00',0,'login_failed','User not found','[::1]:52006','{\"username\": \"puto\"}'),(2,'2024-05-25 15:21:52',0,'register_user_success','User registered successfully','[::1]:52038','{\"username\": \"mamamio\"}'),(3,'2024-05-25 15:21:58',0,'login_failed','User not found','[::1]:52039','{\"username\": \"mamamio\"}'),(4,'2024-05-25 15:33:58',3,'login_failed','Invalid login credentials','[::1]:52117','null'),(5,'2024-05-25 15:34:20',0,'register_user_success','User registered successfully','[::1]:52119','{\"username\": \"hola\"}'),(6,'2024-05-25 15:34:23',4,'login_success','User logged in successfully','[::1]:52121','null'),(7,'2024-05-25 15:34:44',0,'authentication_success','User authenticated successfully','[::1]:52121','{\"username\": \"hola\"}'),(8,'2024-05-25 15:34:44',4,'create_post_success','Post created successfully','[::1]:52121','{\"post_title\": \"guey\"}'),(9,'2024-05-25 15:35:00',0,'authentication_success','User authenticated successfully','[::1]:52156','{\"username\": \"hola\"}'),(10,'2024-05-25 15:35:00',0,'get_posts_success','Posts retrieved successfully','[::1]:52156','null'),(11,'2024-05-25 15:35:06',0,'authentication_success','User authenticated successfully','[::1]:52156','{\"username\": \"hola\"}'),(12,'2024-05-25 15:35:06',0,'get_users_success','Users retrieved successfully','[::1]:52156','null'),(13,'2024-05-25 15:35:07',0,'authentication_success','User authenticated successfully','[::1]:52156','{\"username\": \"hola\"}'),(14,'2024-05-25 15:35:13',0,'authentication_success','User authenticated successfully','[::1]:52156','{\"username\": \"hola\"}'),(15,'2024-05-25 15:35:13',0,'get_users_success','Users retrieved successfully','[::1]:52156','null'),(16,'2024-05-25 15:35:15',0,'authentication_success','User authenticated successfully','[::1]:52156','{\"username\": \"hola\"}'),(17,'2024-05-26 10:17:19',3,'login_failed','Invalid login credentials','[::1]:64208','null'),(18,'2024-05-26 10:17:27',3,'login_failed','Invalid login credentials','[::1]:64209','null'),(19,'2024-05-26 10:17:36',4,'login_success','User logged in successfully','[::1]:64211','null'),(20,'2024-05-26 10:17:38',0,'authentication_success','User authenticated successfully','[::1]:64211','{\"username\": \"hola\"}'),(21,'2024-05-26 10:17:38',0,'get_users_success','Users retrieved successfully','[::1]:64211','null'),(22,'2024-05-26 10:17:39',0,'authentication_success','User authenticated successfully','[::1]:64211','{\"username\": \"hola\"}'),(23,'2024-05-26 10:18:55',0,'register_user_success','User registered successfully','[::1]:64211','{\"username\": \"adri\"}'),(24,'2024-05-26 10:19:10',5,'login_success','User logged in successfully','[::1]:64294','null'),(25,'2024-05-26 12:52:13',5,'login_success','User logged in successfully','[::1]:50942','null'),(26,'2024-05-26 12:53:53',0,'authentication_success','User authenticated successfully','[::1]:50955','{\"username\": \"adri\"}'),(27,'2024-05-26 12:55:50',0,'authentication_success','User authenticated successfully','[::1]:50986','{\"username\": \"adri\"}'),(28,'2024-05-26 13:01:49',0,'authentication_success','User authenticated successfully','[::1]:51047','{\"username\": \"adri\"}'),(29,'2024-05-26 13:02:19',0,'authentication_success','User authenticated successfully','[::1]:51088','{\"username\": \"adri\"}'),(30,'2024-05-26 13:17:43',0,'authentication_success','User authenticated successfully','[::1]:51300','{\"username\": \"adri\"}'),(31,'2024-05-26 13:17:43',0,'backup_failed','mysqldump failed','','{\"error\": \"exec: \\\"mysqldump\\\": executable file not found in %PATH%\"}'),(32,'2024-05-26 13:17:43',0,'backup_failed','Backup failed','[::1]:51300','{\"error\": \"mysqldump failed: exec: \\\"mysqldump\\\": executable file not found in %PATH%\"}'),(33,'2024-05-26 13:31:19',0,'authentication_success','User authenticated successfully','[::1]:51413','{\"username\": \"adri\"}'),(34,'2024-05-26 13:31:19',0,'backup_failed','mysqldump failed','','{\"error\": \"exec: \\\"mysqldump\\\": executable file not found in %PATH%\"}'),(35,'2024-05-26 13:31:19',0,'backup_failed','Backup failed','[::1]:51413','{\"error\": \"mysqldump failed: exec: \\\"mysqldump\\\": executable file not found in %PATH%\"}'),(36,'2024-05-26 14:25:53',5,'login_success','User logged in successfully','[::1]:60225','null'),(37,'2024-05-26 14:25:56',0,'authentication_success','User authenticated successfully','[::1]:60225','{\"username\": \"adri\"}'),(38,'2024-05-27 12:07:58',5,'login_success','User logged in successfully','[::1]:53496','null'),(39,'2024-05-27 12:08:02',0,'authentication_success','User authenticated successfully','[::1]:53496','{\"username\": \"adri\"}');
/*!40000 ALTER TABLE `event_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mensaje`
--

DROP TABLE IF EXISTS `mensaje`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mensaje` (
  `id` int NOT NULL AUTO_INCREMENT,
  `remitente_id` int NOT NULL,
  `destinatario_id` int NOT NULL,
  `contenido` text NOT NULL,
  `fecha_envio` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `remitente_id` (`remitente_id`),
  KEY `destinatario_id` (`destinatario_id`),
  CONSTRAINT `mensaje_ibfk_1` FOREIGN KEY (`remitente_id`) REFERENCES `usuario` (`id`),
  CONSTRAINT `mensaje_ibfk_2` FOREIGN KEY (`destinatario_id`) REFERENCES `usuario` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mensaje`
--

LOCK TABLES `mensaje` WRITE;
/*!40000 ALTER TABLE `mensaje` DISABLE KEYS */;
/*!40000 ALTER TABLE `mensaje` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post`
--

DROP TABLE IF EXISTS `post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post` (
  `id` int NOT NULL AUTO_INCREMENT,
  `titulo` varchar(255) NOT NULL,
  `texto` varchar(255) NOT NULL,
  `usuario_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `usuario_id` (`usuario_id`),
  CONSTRAINT `post_ibfk_1` FOREIGN KEY (`usuario_id`) REFERENCES `usuario` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post`
--

LOCK TABLES `post` WRITE;
/*!40000 ALTER TABLE `post` DISABLE KEYS */;
/*!40000 ALTER TABLE `post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `nombre` (`nombre`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (3,'administrador'),(2,'moderador'),(1,'usuario');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usuario`
--

DROP TABLE IF EXISTS `usuario`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usuario` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `username` varchar(100) NOT NULL,
  `hashed_username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `nombre` varchar(100) DEFAULT NULL,
  `apellidos` varchar(100) DEFAULT NULL,
  `rol_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `rol_id` (`rol_id`),
  CONSTRAINT `usuario_ibfk_1` FOREIGN KEY (`rol_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuario`
--

LOCK TABLES `usuario` WRITE;
/*!40000 ALTER TABLE `usuario` DISABLE KEYS */;
INSERT INTO `usuario` VALUES (5,'35af126dafdfd1567cd9f81195d3b95f6054448b69a9684ec7c95bda8ec6de411a503ca650fe96568456d20a81283b','24b12cc0a49629a4a8cc540f00566212c9e596b72ea7f1ca60f11e74383c21ee','6a044a29433dda71ba45d7ae1f746548ccf7c7978b055a9551aca5eb869e4c86','4dbf558815b21e6678e41715816ca0f1:1cd35894de1060b484d973adde5e20d377c2f09e00a80c58015bcabc8e60e939','6d20975ccdfaf35a1fcd1469da23d949d00a0f1387f345b698c804da415eaa94','3683363259854f8d6b29651775dbc64b4b2158cf19f89c694560f75c2e7d1afa',3);
/*!40000 ALTER TABLE `usuario` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-05-27 14:08:04
