/*
Navicat SQLite Data Transfer

Source Server         : tienda&admin
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2018-01-15 10:34:18
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for aux
-- ----------------------------
DROP TABLE IF EXISTS "main"."aux";
CREATE TABLE "aux" (
"hora_inicial"  INTEGER,
"hora_final"  INTEGER
);

-- ----------------------------
-- Records of aux
-- ----------------------------

-- ----------------------------
-- Table structure for horario
-- ----------------------------
DROP TABLE IF EXISTS "main"."horario";
CREATE TABLE "horario" (
"hora_inicial"  TEXT,
"hora_final"  TEXT
);

-- ----------------------------
-- Records of horario
-- ----------------------------

-- ----------------------------
-- Table structure for mensaje
-- ----------------------------
DROP TABLE IF EXISTS "main"."mensaje";
CREATE TABLE "mensaje" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"fichero"  TEXT(255),
"playtime"  TEXT(5),
"existe"  TEXT(1),
"fecha"  TEXT(10)
);

-- ----------------------------
-- Records of mensaje
-- ----------------------------

-- ----------------------------
-- Table structure for musica
-- ----------------------------
DROP TABLE IF EXISTS "main"."musica";
CREATE TABLE "musica" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"carpetas"  TEXT(4096),
"fecha_inicio"  TEXT(10),
"fecha_final"  TEXT(10),
"timestamp"  INTEGER
);

-- ----------------------------
-- Records of musica
-- ----------------------------

-- ----------------------------
-- Table structure for publi
-- ----------------------------
DROP TABLE IF EXISTS "main"."publi";
CREATE TABLE "publi" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"fichero"  TEXT(255),
"existe"  TEXT(1),
"fecha_ini"  TEXT(10),
"gap"  INTEGER
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
INSERT INTO "main"."sqlite_sequence" VALUES ('musica', 0);
INSERT INTO "main"."sqlite_sequence" VALUES ('mensaje', 0);
INSERT INTO "main"."sqlite_sequence" VALUES ('publi', 0);
INSERT INTO "main"."sqlite_sequence" VALUES ('tienda', 0);
INSERT INTO "main"."sqlite_sequence" VALUES ('usuarios', 0);

-- ----------------------------
-- Table structure for tienda
-- ----------------------------
DROP TABLE IF EXISTS "main"."tienda";
CREATE TABLE "tienda" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"dominio"  TEXT,
"last_connect"  INTEGER
);

-- ----------------------------
-- Records of tienda
-- ----------------------------

-- ----------------------------
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS "main"."usuarios";
CREATE TABLE "usuarios" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"user"  TEXT(25),
"pass"  TEXT(25)
);

-- ----------------------------
-- Records of usuarios
-- ----------------------------
