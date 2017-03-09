/*
Navicat SQLite Data Transfer

Source Server         : Prueba
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2017-03-09 16:29:33
*/

PRAGMA foreign_keys = OFF;

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
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);
