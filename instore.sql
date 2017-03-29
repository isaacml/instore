/*
Navicat SQLite Data Transfer

Source Server         : INSTORE
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2017-03-29 12:11:04
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for almacenes
-- ----------------------------
DROP TABLE IF EXISTS "main"."almacenes";
CREATE TABLE "almacenes" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"almacen"  TEXT(255),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"entidad_id"  INTEGER NOT NULL,
CONSTRAINT "fk_usuario" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT "fk_entidad" FOREIGN KEY ("entidad_id") REFERENCES "entidades" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for entidades
-- ----------------------------
DROP TABLE IF EXISTS "main"."entidades";
CREATE TABLE "entidades" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"nombre"  TEXT(255),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"last_access"  INTEGER,
CONSTRAINT "fk_ent_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for mensaje
-- ----------------------------
DROP TABLE IF EXISTS "main"."mensaje";
CREATE TABLE "mensaje" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"fichero"  TEXT(255),
"fecha_inicio"  TEXT(10),
"fecha_final"  TEXT(10),
"destino"  TEXT(1024),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"playtime"  TEXT(5),
CONSTRAINT "fk_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for pais
-- ----------------------------
DROP TABLE IF EXISTS "main"."pais";
CREATE TABLE "pais" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"pais"  TEXT(255),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"almacen_id"  INTEGER NOT NULL,
CONSTRAINT "fk_usuario" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT "fk_almacen" FOREIGN KEY ("almacen_id") REFERENCES "almacenes" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for provincia
-- ----------------------------
DROP TABLE IF EXISTS "main"."provincia";
CREATE TABLE "provincia" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"provincia"  TEXT(255),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"region_id"  INTEGER NOT NULL,
CONSTRAINT "fk_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT "fk_region" FOREIGN KEY ("region_id") REFERENCES "region" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for publi
-- ----------------------------
DROP TABLE IF EXISTS "main"."publi";
CREATE TABLE "publi" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"fichero"  TEXT(255),
"fecha_inicio"  TEXT(10),
"fecha_final"  TEXT(10),
"destino"  TEXT(1024),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"gap"  INTEGER,
CONSTRAINT "fk_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for region
-- ----------------------------
DROP TABLE IF EXISTS "main"."region";
CREATE TABLE "region" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"region"  TEXT(255),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"pais_id"  INTEGER NOT NULL,
CONSTRAINT "fk_usuario" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT "fk_pais" FOREIGN KEY ("pais_id") REFERENCES "pais" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);

-- ----------------------------
-- Table structure for tiendas
-- ----------------------------
DROP TABLE IF EXISTS "main"."tiendas";
CREATE TABLE "tiendas" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"tienda"  TEXT(255),
"creador_id"  INTEGER NOT NULL,
"timestamp"  INTEGER,
"provincia_id"  INTEGER NOT NULL,
"address"  TEXT(1024),
"phone"  TEXT(255),
"extra"  TEXT(1024),
CONSTRAINT "fk_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT "fk_prov" FOREIGN KEY ("provincia_id") REFERENCES "provincia" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS "main"."usuarios";
CREATE TABLE "usuarios" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"user"  TEXT(25),
"old_user"  TEXT(25),
"pass"  TEXT(25),
"nombre_completo"  TEXT(255),
"entidad_id"  INTEGER,
"padre_id"  INTEGER,
"bitmap_acciones"  TEXT(16),
CONSTRAINT "fk_user_entidad" FOREIGN KEY ("entidad_id") REFERENCES "entidades" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);
