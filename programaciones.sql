/*
Navicat SQLite Data Transfer

Source Server         : Programaciones
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2017-03-13 21:03:19
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for mensaje
-- ----------------------------
DROP TABLE IF EXISTS "main"."mensaje";
CREATE TABLE "mensaje" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"ruta"  TEXT(1024),
"fichero"  TEXT(255),
"fecha_inicio"  TEXT(10),
"fecha_final"  TEXT(10),
"destino"  TEXT(1024),
"timestamp"  INTEGER,
"playtime"  TEXT(5)
);

-- ----------------------------
-- Records of mensaje
-- ----------------------------

-- ----------------------------
-- Table structure for publi
-- ----------------------------
DROP TABLE IF EXISTS "main"."publi";
CREATE TABLE "publi" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"ruta"  TEXT(1024),
"fichero"  TEXT(255),
"fecha_inicio"  TEXT(10),
"fecha_final"  TEXT(10),
"destino"  TEXT(1024),
"timestamp"  INTEGER
);

-- ----------------------------
-- Records of publi
-- ----------------------------

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);

-- ----------------------------
-- Records of sqlite_sequence
-- ----------------------------
INSERT INTO "main"."sqlite_sequence" VALUES ('publi', 1);
INSERT INTO "main"."sqlite_sequence" VALUES ('mensaje', 0);
