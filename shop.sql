/*
Navicat SQLite Data Transfer

Source Server         : shop mejorada
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2018-03-01 06:20:34
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
"existe"  TEXT(1),
"fecha_ini"  TEXT(10),
"fecha_fin"  TEXT(10),
"playtime"  TEXT(5)
);

-- ----------------------------
-- Table structure for musica
-- ----------------------------
DROP TABLE IF EXISTS "main"."musica";
CREATE TABLE "musica" (
"carpeta"  TEXT(4096)
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
"fecha_fin"  TEXT(10),
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
