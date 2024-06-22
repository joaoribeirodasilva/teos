/*!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.6.18-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: teos
-- ------------------------------------------------------
-- Server version	10.6.18-MariaDB-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `app_configurations`
--

DROP TABLE IF EXISTS `app_configurations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_configurations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_environment_id` bigint(20) unsigned NOT NULL,
  `application_id` bigint(20) unsigned NOT NULL,
  `configuration_key` varchar(100) NOT NULL,
  `type` enum('string','int','double','bool','date','time','datetime') NOT NULL,
  `val_string` text DEFAULT NULL,
  `val_int` bigint(20) DEFAULT NULL,
  `val_double` double DEFAULT NULL,
  `val_boolean` tinyint(1) DEFAULT 0,
  `val_date` date DEFAULT NULL,
  `val_time` time DEFAULT NULL,
  `val_datetime` timestamp NULL DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `app_configurations_app_environment_id_IDX` (`app_environment_id`,`application_id`,`configuration_key`) USING BTREE,
  KEY `app_configurations_applications_FK` (`application_id`),
  KEY `app_configurations_deletedAt_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `app_configurations_app_environments_FK` FOREIGN KEY (`app_environment_id`) REFERENCES `app_environments` (`id`),
  CONSTRAINT `app_configurations_applications_FK` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_configurations`
--

LOCK TABLES `app_configurations` WRITE;
/*!40000 ALTER TABLE `app_configurations` DISABLE KEYS */;
INSERT INTO `app_configurations` VALUES (2,2,1,'AUTH_COOKIE_NAME','string','teos_auth',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:39:15',NULL,NULL),(3,2,1,'AUTH_COOKIE_EXPIRE','int',NULL,900,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:39:15',NULL,NULL),(4,2,1,'AUTH_COOKIE_HTTP_ONLY','bool',NULL,NULL,NULL,1,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:39:15',NULL,NULL),(5,2,1,'AUTH_COOKIE_SECURE','bool',NULL,NULL,NULL,1,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:39:15',NULL,NULL),(6,2,1,'AUTH_SECRET_KEY','string','bwfeb778(*!$dag#Fahjd',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:39:15',NULL,NULL),(7,2,2,'HTTP_BIND_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL),(8,2,2,'HTTP_BIND_PORT','int',NULL,8080,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL),(9,2,3,'HTTP_BIND_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL),(10,2,3,'HTTP_BIND_PORT','int',NULL,8081,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL),(11,2,4,'HTTP_BIND_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL),(12,2,4,'HTTP_BIND_PORT','int',NULL,8082,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL),(13,2,5,'HTTP_BIND_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(14,2,5,'HTTP_BIND_PORT','int',NULL,8083,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(15,2,1,'DB_PERMISSIONS_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(16,2,1,'DB_PERMISSIONS_PORT','int',NULL,6379,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(17,2,1,'DB_PERMISSIONS_DB','int',NULL,0,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(18,2,1,'DB_SESSIONS_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(19,2,1,'DB_SESSIONS_PORT','int',NULL,6379,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(20,2,1,'DB_SESSIONS_DB','int',NULL,1,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(21,2,1,'DB_HISTORY_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(22,2,1,'DB_HISTORY_PORT','int',NULL,6379,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(23,2,1,'DB_HISTORY_DB','int',NULL,2,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(24,2,1,'DB_LOGS_ADDR','string','localhost',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(25,2,1,'DB_LOGS_PORT','int',NULL,6379,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(26,2,1,'DB_LOGS_DB','int',NULL,3,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 18:00:07',NULL,NULL),(27,2,1,'AUTH_COOKIE_DOMAIN','string','',NULL,NULL,0,NULL,NULL,NULL,1,'2024-06-21 18:33:20',1,'2024-06-19 17:54:50',NULL,NULL);
/*!40000 ALTER TABLE `app_configurations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_environments`
--

DROP TABLE IF EXISTS `app_environments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_environments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `key` varchar(50) NOT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `app_environments_name_IDX` (`name`) USING BTREE,
  KEY `app_environments_deletedAt_IDX` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_environments`
--

LOCK TABLES `app_environments` WRITE;
/*!40000 ALTER TABLE `app_environments` DISABLE KEYS */;
INSERT INTO `app_environments` VALUES (1,'Desenvolvimento','Ambiente de desenvolvimento','dev',1,'2024-06-19 17:39:00',1,'2024-06-19 17:22:11',NULL,NULL),(2,'Dev Local','Ambiente de desenvolviment local','local',1,'2024-06-21 18:35:06',1,'2024-06-19 17:39:00',NULL,NULL),(3,'Teste','Ambiente ded testes','test',1,'2024-06-21 18:30:04',1,'2024-06-19 17:22:11',NULL,NULL),(4,'Avaliação','Ambiente de avaliação','stage',1,'2024-06-21 18:29:59',1,'2024-06-19 17:22:11',NULL,NULL),(5,'Produção','Ambiente de produção','prod',1,'2024-06-21 18:29:55',1,'2024-06-19 17:22:11',NULL,NULL);
/*!40000 ALTER TABLE `app_environments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_route_methods`
--

DROP TABLE IF EXISTS `app_route_methods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_route_methods` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_route_id` bigint(20) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `method` enum('GET','POST','PUT','PATCH','DELETE') NOT NULL DEFAULT 'GET',
  `uri` varchar(255) NOT NULL,
  `open` tinyint(1) NOT NULL DEFAULT 0,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_route_methods_unique_IDX` (`app_route_id`,`method`,`uri`) USING BTREE,
  UNIQUE KEY `app_route_methods_deletedAt_IDX` (`deleted_at`) USING BTREE,
  KEY `app_route_methods_app_route_id_IDX` (`app_route_id`,`name`) USING BTREE,
  CONSTRAINT `app_route_methods_app_routes_FK` FOREIGN KEY (`app_route_id`) REFERENCES `app_routes` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_route_methods`
--

LOCK TABLES `app_route_methods` WRITE;
/*!40000 ALTER TABLE `app_route_methods` DISABLE KEYS */;
INSERT INTO `app_route_methods` VALUES (1,1,'Listar','Lista todos os registros de dados','GET','',0,1,1,'2024-06-20 18:38:54',1,'2024-06-20 18:20:38',NULL,NULL),(2,1,'Visualizar','Visualiza os detalhes dos registros de dados','GET','/:id',0,1,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(3,1,'Criar','Cria novos registros de dados','POST','',0,1,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(4,1,'Alterar','Altera registros de dados existentes','PUT','/:id',0,1,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(5,1,'Apagar','Apaga registros de dados','DELETE','/:id',0,1,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(6,2,'Listar','Lista todos os registros de dados','GET','',0,0,1,'2024-06-20 18:38:54',1,'2024-06-20 18:20:38',NULL,NULL),(7,2,'Visualizar','Visualiza os detalhes dos registros de dados','GET','/:id',0,0,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(8,2,'Criar','Cria novos registros de dados','POST','',0,0,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(9,2,'Alterar','Altera registros de dados existentes','PUT','/:id',0,0,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(10,2,'Apagar','Apaga registros de dados','DELETE','/:id',0,0,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(13,3,'Listar','Lista todos os registros de dados','GET','',0,0,1,'2024-06-20 18:38:54',1,'2024-06-20 18:20:38',NULL,NULL),(14,3,'Visualizar','Visualiza os detalhes dos registros de dados','GET','/:id',0,0,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(15,3,'Criar','Cria novos registros de dados','POST','',0,0,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(16,3,'Alterar','Altera registros de dados existentes','PUT','/:id',0,0,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(17,3,'Apagar','Apaga registros de dados','DELETE','/:id',0,0,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(20,4,'Listar','Lista todos os registros de dados','GET','',0,0,1,'2024-06-20 18:38:54',1,'2024-06-20 18:20:38',NULL,NULL),(21,4,'Visualizar','Visualiza os detalhes dos registros de dados','GET','/:id',0,0,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(22,4,'Criar','Cria novos registros de dados','POST','',0,0,1,'2024-06-20 18:38:44',1,'2024-06-20 18:20:38',NULL,NULL),(23,4,'Alterar','Altera registros de dados existentes','PUT','/:id',0,0,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(24,4,'Apagar','Apaga registros de dados','DELETE','/:id',0,0,1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL);
/*!40000 ALTER TABLE `app_route_methods` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_routes`
--

DROP TABLE IF EXISTS `app_routes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_routes` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `application_id` bigint(20) unsigned NOT NULL,
  `app_environment_id` bigint(20) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `uri` varchar(255) NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_routes_application_id_IDX` (`application_id`,`app_environment_id`,`name`) USING BTREE,
  UNIQUE KEY `app_routes_unique2_IDX` (`application_id`,`app_environment_id`,`uri`) USING BTREE,
  KEY `app_routes_app_environments_FK` (`app_environment_id`),
  KEY `app_routes_deletedAt_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `app_routes_app_configurations_FK` FOREIGN KEY (`application_id`) REFERENCES `app_configurations` (`id`),
  CONSTRAINT `app_routes_app_environments_FK` FOREIGN KEY (`app_environment_id`) REFERENCES `app_environments` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_routes`
--

LOCK TABLES `app_routes` WRITE;
/*!40000 ALTER TABLE `app_routes` DISABLE KEYS */;
INSERT INTO `app_routes` VALUES (1,4,2,'Usuários','Cadastro de usuários','/users/users',1,1,'2024-06-21 18:34:06',1,'2024-06-20 21:21:00',NULL,NULL),(2,4,2,'Grupos','Grupos dos usuários','/users/groups',1,1,'2024-06-21 18:34:06',1,'2024-06-20 21:21:00',NULL,NULL),(3,4,2,'Organizações','Organizações dos usuários','/users/organizations',1,1,'2024-06-21 18:34:06',1,'2024-06-20 18:20:38',NULL,NULL),(4,4,2,'Permissões','Permissões de acesso específicas a usuários','/user/permissions',1,1,'2024-06-21 18:34:06',1,'2024-06-20 18:20:38',NULL,NULL);
/*!40000 ALTER TABLE `app_routes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `applications`
--

DROP TABLE IF EXISTS `applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `applications` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `code` varchar(255) NOT NULL,
  `internal` tinyint(1) NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `applications_deletedAt_IDX` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applications`
--

LOCK TABLES `applications` WRITE;
/*!40000 ALTER TABLE `applications` DISABLE KEYS */;
INSERT INTO `applications` VALUES (1,'[GLOBAL]','Todos os serviços','global',1,1,'2024-06-19 17:39:56',1,'2024-06-19 17:39:15',NULL,NULL),(2,'Autenticação','Serviço de autenticação e autorização de usuários','auth',0,1,'2024-06-19 17:40:00',1,'2024-06-19 21:16:00',NULL,NULL),(3,'Histórico','Serviço de histórico de alteração de dados','hist',0,1,'2024-06-19 17:40:03',1,'2024-06-19 21:16:00',NULL,NULL),(4,'Usuários','Serviço de dados de usuários','user',0,1,'2024-06-19 17:40:07',1,'2024-06-19 21:16:00',NULL,NULL),(5,'Applicações','Serviço de dados de aplicações','apps',1,1,'2024-06-19 21:35:20',1,'2024-06-19 21:16:00',NULL,NULL),(6,'Consumidor de históricos','Serviço que lê dados de histórico de alterações do redis e os insere no banco','chists',1,1,'2024-06-22 05:47:37',1,'2024-06-19 17:39:56',NULL,NULL),(7,'Consumidor de logs','Serviço que lê dados de logs e os insere no banco','clogs',1,1,'2024-06-22 05:51:21',1,'2024-06-22 05:47:37',NULL,NULL),(8,'Consumidor de permissões','Serviço que consome permissões de usuários e os disponibiliza no redis','cperms',1,1,'2024-06-22 05:51:21',1,'2024-06-22 05:47:37',NULL,NULL);
/*!40000 ALTER TABLE `applications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_groups`
--

DROP TABLE IF EXISTS `auth_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_groups` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_groups_organization_id_IDX` (`organization_id`,`name`) USING BTREE,
  KEY `auth_groups_deletedAt_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `auth_groups_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_groups`
--

LOCK TABLES `auth_groups` WRITE;
/*!40000 ALTER TABLE `auth_groups` DISABLE KEYS */;
INSERT INTO `auth_groups` VALUES (1,1,'Global','Grupo de acesso global, acessa todas as áreas da aplicação',1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL),(2,2,'Global','Grupo de acesso global, acessa todas as áreas da aplicação',1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL);
/*!40000 ALTER TABLE `auth_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_roles`
--

DROP TABLE IF EXISTS `auth_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `roles_organization_id_IDX` (`organization_id`,`name`) USING BTREE,
  KEY `auth_roles_deletedAt_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `roles_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_roles`
--

LOCK TABLES `auth_roles` WRITE;
/*!40000 ALTER TABLE `auth_roles` DISABLE KEYS */;
INSERT INTO `auth_roles` VALUES (1,2,'Administrador','Administrador geral do sistema',1,'2024-06-20 18:24:43',1,'2024-06-20 18:20:38',NULL,NULL);
/*!40000 ALTER TABLE `auth_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_sessions`
--

DROP TABLE IF EXISTS `auth_sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_sessions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `sessions_users_FK` (`user_id`),
  KEY `auth_sessions_deletedAt_IDX` (`deleted_at`) USING BTREE,
  KEY `auth_sessions_organizations_FK` (`organization_id`),
  CONSTRAINT `auth_sessions_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`),
  CONSTRAINT `sessions_users_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_sessions`
--

LOCK TABLES `auth_sessions` WRITE;
/*!40000 ALTER TABLE `auth_sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `history`
--

DROP TABLE IF EXISTS `history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `history` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `table_name` varchar(255) NOT NULL,
  `original_id` bigint(20) unsigned NOT NULL,
  `json_data` text NOT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `history_deletedAt_IDX` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `history`
--

LOCK TABLES `history` WRITE;
/*!40000 ALTER TABLE `history` DISABLE KEYS */;
/*!40000 ALTER TABLE `history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `org_applications`
--

DROP TABLE IF EXISTS `org_applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `org_applications` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `application_id` bigint(20) unsigned NOT NULL,
  `app_environment_id` bigint(20) unsigned NOT NULL,
  `active` tinyint(4) NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `org_applications_organization_id_IDX` (`organization_id`,`application_id`,`app_environment_id`) USING BTREE,
  KEY `org_applications_applications_FK` (`application_id`),
  KEY `org_applications_app_environments_FK` (`app_environment_id`),
  KEY `org_applications_deletedAt_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `org_applications_app_environments_FK` FOREIGN KEY (`app_environment_id`) REFERENCES `app_environments` (`id`),
  CONSTRAINT `org_applications_applications_FK` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`),
  CONSTRAINT `org_applications_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `org_applications`
--

LOCK TABLES `org_applications` WRITE;
/*!40000 ALTER TABLE `org_applications` DISABLE KEYS */;
/*!40000 ALTER TABLE `org_applications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `org_types`
--

DROP TABLE IF EXISTS `org_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `org_types` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `organization_types_organization_id_IDX` (`organization_id`,`name`) USING BTREE,
  KEY `org_types_deleted_At_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `organization_types_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `org_types`
--

LOCK TABLES `org_types` WRITE;
/*!40000 ALTER TABLE `org_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `org_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organizations`
--

DROP TABLE IF EXISTS `organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `organizations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organizations`
--

LOCK TABLES `organizations` WRITE;
/*!40000 ALTER TABLE `organizations` DISABLE KEYS */;
INSERT INTO `organizations` VALUES (1,'Teos',1,'2024-06-19 17:18:40',1,'2024-06-19 21:18:00',NULL,NULL),(2,'Universidade Teos',1,'2024-06-20 18:20:38',1,'2024-06-20 18:20:38',NULL,NULL);
/*!40000 ALTER TABLE `organizations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reset_types`
--

DROP TABLE IF EXISTS `reset_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `reset_types` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `reset_types_deleted_At_IDX` (`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reset_types`
--

LOCK TABLES `reset_types` WRITE;
/*!40000 ALTER TABLE `reset_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `reset_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `resets`
--

DROP TABLE IF EXISTS `resets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resets` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `reset_type_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `code` varchar(255) NOT NULL,
  `expire` int(10) unsigned NOT NULL,
  `used` timestamp NULL DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `resets_users_FK` (`user_id`),
  KEY `resets_reset_types_FK` (`reset_type_id`),
  KEY `resets_deleted_At_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `resets_reset_types_FK` FOREIGN KEY (`reset_type_id`) REFERENCES `reset_types` (`id`),
  CONSTRAINT `resets_users_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `resets`
--

LOCK TABLES `resets` WRITE;
/*!40000 ALTER TABLE `resets` DISABLE KEYS */;
/*!40000 ALTER TABLE `resets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_groups`
--

DROP TABLE IF EXISTS `user_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_groups` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `auth_group_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_groups_auth_organization_id_IDX` (`organization_id`,`auth_group_id`,`user_id`) USING BTREE,
  KEY `user_groups_auth_groups_FK` (`auth_group_id`),
  KEY `user_groups_users_FK` (`user_id`),
  KEY `user_groups_deleted_At_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `user_groups_auth_groups_FK` FOREIGN KEY (`auth_group_id`) REFERENCES `auth_groups` (`id`),
  CONSTRAINT `user_groups_org_applications_FK` FOREIGN KEY (`organization_id`) REFERENCES `org_applications` (`id`),
  CONSTRAINT `user_groups_users_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_groups`
--

LOCK TABLES `user_groups` WRITE;
/*!40000 ALTER TABLE `user_groups` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_organizations`
--

DROP TABLE IF EXISTS `user_organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_organizations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `active` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_organizations_organization_id_IDX` (`organization_id`,`user_id`) USING BTREE,
  KEY `user_organizations_users_FK` (`user_id`),
  KEY `user_organizations_deleted_At_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `user_organizations_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`),
  CONSTRAINT `user_organizations_users_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_organizations`
--

LOCK TABLES `user_organizations` WRITE;
/*!40000 ALTER TABLE `user_organizations` DISABLE KEYS */;
INSERT INTO `user_organizations` VALUES (1,1,1,1,1,'2024-06-19 17:19:19',1,'2024-06-19 21:19:00',NULL,NULL);
/*!40000 ALTER TABLE `user_organizations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_permissions`
--

DROP TABLE IF EXISTS `user_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` bigint(20) unsigned NOT NULL,
  `auth_role_id` bigint(20) unsigned NOT NULL,
  `app_route_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT 0,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_permissions_app_routes_FK` (`app_route_id`),
  KEY `user_permissions_users_FK` (`user_id`),
  KEY `user_permissions_auth_roles_FK` (`auth_role_id`),
  KEY `user_permissions_organization_id_IDX` (`organization_id`,`auth_role_id`,`app_route_id`,`user_id`) USING BTREE,
  KEY `user_permissions_deleted_At_IDX` (`deleted_at`) USING BTREE,
  CONSTRAINT `user_permissions_app_routes_FK` FOREIGN KEY (`app_route_id`) REFERENCES `app_routes` (`id`),
  CONSTRAINT `user_permissions_auth_roles_FK` FOREIGN KEY (`auth_role_id`) REFERENCES `auth_roles` (`id`),
  CONSTRAINT `user_permissions_organizations_FK` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`),
  CONSTRAINT `user_permissions_users_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_permissions`
--

LOCK TABLES `user_permissions` WRITE;
/*!40000 ALTER TABLE `user_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) NOT NULL,
  `surname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) DEFAULT NULL,
  `terms` timestamp NULL DEFAULT NULL,
  `avatar_url` varchar(255) DEFAULT '',
  `email_verified` timestamp NULL DEFAULT NULL,
  `created_by` bigint(20) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `updated_by` bigint(20) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_by` bigint(20) unsigned DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_IDX` (`email`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'Sistema','Teos','teos@teos.com.br',' ',NULL,'',NULL,1,'2024-06-19 23:39:18',1,'2024-06-19 21:17:00',NULL,NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'teos'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-06-22  2:56:31
