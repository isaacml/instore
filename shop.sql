/*
Navicat SQLite Data Transfer

Source Server         : Programaciones
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2017-12-25 21:02:50
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for horario
-- ----------------------------
DROP TABLE IF EXISTS "main"."horario";
CREATE TABLE "horario" (
"hora_inicial"  TEXT,
"hora_final"  TEXT
);

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
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);

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
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS "main"."usuarios";
CREATE TABLE "usuarios" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"user"  TEXT(25),
"pass"  TEXT(25)
);
