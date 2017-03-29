/*
Navicat SQLite Data Transfer

Source Server         : Musica
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2017-03-29 12:10:20
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for musica
-- ----------------------------
DROP TABLE IF EXISTS "main"."musica";
CREATE TABLE "musica" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"carpetas"  TEXT(4096),
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
