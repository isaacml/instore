/*
Navicat SQLite Data Transfer

Source Server         : Servidor Interno
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2018-01-15 10:26:10
*/

PRAGMA foreign_keys = OFF;

-- ----------------------------
-- Table structure for mensaje
-- ----------------------------
DROP TABLE IF EXISTS "main"."mensaje";
CREATE TABLE "mensaje" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"fichero"  TEXT(255),
"existe"  TEXT(1),
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
"fecha_ini"  TEXT,
"timestamp"  INTEGER,
"gap"  INTEGER
);

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);

-- ----------------------------
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS "main"."usuarios";
CREATE TABLE "usuarios" (
"id"  INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"user"  TEXT,
"pass"  TEXT
);
