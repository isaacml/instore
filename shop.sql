/*
 Navicat SQLite Data Transfer

 Source Server         : pruebas JR2
 Source Server Type    : SQLite
 Source Server Version : 3021000
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3021000
 File Encoding         : 65001

 Date: 12/03/2018 12:33:42
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for aux
-- ----------------------------
DROP TABLE IF EXISTS "aux";
CREATE TABLE "aux" (
  "hora_inicial" INTEGER,
  "hora_final" INTEGER
);

-- ----------------------------
-- Table structure for horario
-- ----------------------------
DROP TABLE IF EXISTS "horario";
CREATE TABLE "horario" (
  "hora_inicial" TEXT,
  "hora_final" TEXT
);

-- ----------------------------
-- Table structure for mensaje
-- ----------------------------
DROP TABLE IF EXISTS "mensaje";
CREATE TABLE "mensaje" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "fichero" TEXT(255),
  "existe" TEXT(1),
  "fecha_ini" TEXT(10),
  "fecha_fin" TEXT(10),
  "playtime" TEXT(5)
);

-- ----------------------------
-- Table structure for musica
-- ----------------------------
DROP TABLE IF EXISTS "musica";
CREATE TABLE "musica" (
  "carpeta" TEXT(4096)
);

-- ----------------------------
-- Table structure for publi
-- ----------------------------
DROP TABLE IF EXISTS "publi";
CREATE TABLE "publi" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "fichero" TEXT(255),
  "existe" TEXT(1),
  "fecha_ini" TEXT(10),
  "fecha_fin" TEXT(10),
  "gap" INTEGER
);

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "sqlite_sequence";
CREATE TABLE "sqlite_sequence" (
  "name",
  "seq"
);

-- ----------------------------
-- Records of sqlite_sequence
-- ----------------------------
INSERT INTO "sqlite_sequence" VALUES ('tienda', 0);
INSERT INTO "sqlite_sequence" VALUES ('usuarios', 0);
INSERT INTO "sqlite_sequence" VALUES ('publi', 0);
INSERT INTO "sqlite_sequence" VALUES ('mensaje', 0);

-- ----------------------------
-- Table structure for st_prog_music
-- ----------------------------
DROP TABLE IF EXISTS "st_prog_music";
CREATE TABLE "st_prog_music" (
  "estado" TEXT
);

-- ----------------------------
-- Table structure for tienda
-- ----------------------------
DROP TABLE IF EXISTS "tienda";
CREATE TABLE "tienda" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "dominio" TEXT,
  "last_connect" INTEGER
);

-- ----------------------------
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS "usuarios";
CREATE TABLE "usuarios" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "user" TEXT(25),
  "pass" TEXT(25)
);

-- ----------------------------
-- Auto increment value for mensaje
-- ----------------------------

-- ----------------------------
-- Auto increment value for publi
-- ----------------------------

-- ----------------------------
-- Auto increment value for tienda
-- ----------------------------

-- ----------------------------
-- Auto increment value for usuarios
-- ----------------------------

PRAGMA foreign_keys = true;
