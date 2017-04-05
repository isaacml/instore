/*
Navicat SQLite Data Transfer

Source Server         : TIENDA
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2017-04-05 19:25:23
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for mensaje
-- ----------------------------
DROP TABLE IF EXISTS "main"."mensaje";
CREATE TABLE "mensaje" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"fichero"  TEXT(255),
"playtime"  TEXT(5),
"estado"  TEXT(1)
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
"existe"  TEXT(1)
);

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);
