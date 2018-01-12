/*
Navicat SQLite Data Transfer

Source Server         : INSTORE
Source Server Version : 30808
Source Host           : :0

Target Server Type    : SQLite
Target Server Version : 30808
File Encoding         : 65001

Date: 2018-01-11 03:52:15
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
-- Records of almacenes
-- ----------------------------
INSERT INTO "main"."almacenes" VALUES (1, 'Supersol', 2, 1515409888, 2);
INSERT INTO "main"."almacenes" VALUES (2, 'Cash Diplo', 2, 1515409897, 2);
INSERT INTO "main"."almacenes" VALUES (3, 'Transmediterranea', 2, 1515409913, 3);

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
"status"  INTEGER,
CONSTRAINT "fk_ent_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Records of entidades
-- ----------------------------
INSERT INTO "main"."entidades" VALUES (1, 'Pymedia', 1, 1515409816, 1515409816, 1);
INSERT INTO "main"."entidades" VALUES (2, 'Dinosol', 2, 1515409857, 1515409857, 1);
INSERT INTO "main"."entidades" VALUES (3, 'Acciona', 2, 1515409865, 1515409865, 1);

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
-- Records of mensaje
-- ----------------------------

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
-- Records of pais
-- ----------------------------
INSERT INTO "main"."pais" VALUES (1, 'España', 2, 1515409935, 1);
INSERT INTO "main"."pais" VALUES (2, 'España', 2, 1515409941, 2);
INSERT INTO "main"."pais" VALUES (3, 'España', 2, 1515409946, 3);

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
-- Records of provincia
-- ----------------------------
INSERT INTO "main"."provincia" VALUES (1, 'Almeria', 2, 1515411010, 1);
INSERT INTO "main"."provincia" VALUES (2, 'Cadiz', 2, 1515411021, 1);
INSERT INTO "main"."provincia" VALUES (3, 'Cordoba', 2, 1515411032, 1);
INSERT INTO "main"."provincia" VALUES (4, 'Granada', 2, 1515411042, 1);
INSERT INTO "main"."provincia" VALUES (5, 'Huelva', 2, 1515411085, 1);
INSERT INTO "main"."provincia" VALUES (6, 'Jaen', 2, 1515411101, 1);
INSERT INTO "main"."provincia" VALUES (7, 'Malaga', 2, 1515411111, 1);
INSERT INTO "main"."provincia" VALUES (8, 'Sevilla', 2, 1515411123, 1);
INSERT INTO "main"."provincia" VALUES (9, 'Huesca', 2, 1515411169, 2);
INSERT INTO "main"."provincia" VALUES (10, 'Teruel', 2, 1515411180, 2);
INSERT INTO "main"."provincia" VALUES (11, 'Zaragoza', 2, 1515411189, 2);
INSERT INTO "main"."provincia" VALUES (12, 'Asturias', 2, 1515411262, 3);
INSERT INTO "main"."provincia" VALUES (13, 'Palma de Mallorca', 2, 1515411307, 4);
INSERT INTO "main"."provincia" VALUES (14, 'Las Palmas', 2, 1515411329, 5);
INSERT INTO "main"."provincia" VALUES (15, 'Tenerife', 2, 1515411342, 5);
INSERT INTO "main"."provincia" VALUES (16, 'Cantabria', 2, 1515411378, 6);
INSERT INTO "main"."provincia" VALUES (17, 'Albacete', 2, 1515411396, 7);
INSERT INTO "main"."provincia" VALUES (18, 'Ciudad Real', 2, 1515411413, 7);
INSERT INTO "main"."provincia" VALUES (19, 'Cuenca', 2, 1515411424, 7);
INSERT INTO "main"."provincia" VALUES (20, 'Guadalajara', 2, 1515411442, 7);
INSERT INTO "main"."provincia" VALUES (21, 'Toledo', 2, 1515411451, 7);
INSERT INTO "main"."provincia" VALUES (22, 'Avila', 2, 1515411478, 8);
INSERT INTO "main"."provincia" VALUES (23, 'Burgos', 2, 1515411488, 8);
INSERT INTO "main"."provincia" VALUES (24, 'Leon', 2, 1515411500, 8);
INSERT INTO "main"."provincia" VALUES (25, 'Palencia', 2, 1515411517, 8);
INSERT INTO "main"."provincia" VALUES (26, 'Salamanca', 2, 1515411530, 8);
INSERT INTO "main"."provincia" VALUES (27, 'Segovia', 2, 1515411540, 8);
INSERT INTO "main"."provincia" VALUES (28, 'Soria', 2, 1515411550, 8);
INSERT INTO "main"."provincia" VALUES (29, 'Valladolid', 2, 1515411577, 8);
INSERT INTO "main"."provincia" VALUES (30, 'Zamora', 2, 1515411591, 8);
INSERT INTO "main"."provincia" VALUES (31, 'Barcelona', 2, 1515411603, 9);
INSERT INTO "main"."provincia" VALUES (32, 'Gerona', 2, 1515411620, 9);
INSERT INTO "main"."provincia" VALUES (33, 'Lerida', 2, 1515411630, 9);
INSERT INTO "main"."provincia" VALUES (34, 'Tarragona', 2, 1515411639, 9);
INSERT INTO "main"."provincia" VALUES (35, 'Alicante', 2, 1515411653, 10);
INSERT INTO "main"."provincia" VALUES (36, 'Castellon', 2, 1515411664, 10);
INSERT INTO "main"."provincia" VALUES (37, 'Valencia', 2, 1515411672, 10);
INSERT INTO "main"."provincia" VALUES (38, 'Badajoz', 2, 1515411686, 11);
INSERT INTO "main"."provincia" VALUES (39, 'Caceres', 2, 1515411699, 11);
INSERT INTO "main"."provincia" VALUES (40, 'La Coruña', 2, 1515411714, 12);
INSERT INTO "main"."provincia" VALUES (41, 'Lugo', 2, 1515411733, 12);
INSERT INTO "main"."provincia" VALUES (42, 'Orense', 2, 1515411751, 12);
INSERT INTO "main"."provincia" VALUES (43, 'Pontevedra', 2, 1515411766, 12);
INSERT INTO "main"."provincia" VALUES (44, 'La Rioja', 2, 1515411787, 13);
INSERT INTO "main"."provincia" VALUES (45, 'Madrid', 2, 1515411800, 14);
INSERT INTO "main"."provincia" VALUES (46, 'Murcia', 2, 1515411809, 15);
INSERT INTO "main"."provincia" VALUES (47, 'Navarra', 2, 1515411826, 16);
INSERT INTO "main"."provincia" VALUES (48, 'Alava', 2, 1515411841, 17);
INSERT INTO "main"."provincia" VALUES (49, 'Guipuzcoa', 2, 1515411860, 17);
INSERT INTO "main"."provincia" VALUES (50, 'Vizcaya', 2, 1515411875, 17);
INSERT INTO "main"."provincia" VALUES (51, 'Ceuta', 2, 1515411890, 18);
INSERT INTO "main"."provincia" VALUES (52, 'Melilla', 2, 1515411897, 19);
INSERT INTO "main"."provincia" VALUES (53, 'Almeria', 2, 1515411942, 20);
INSERT INTO "main"."provincia" VALUES (54, 'Cadiz', 2, 1515411955, 20);
INSERT INTO "main"."provincia" VALUES (55, 'Cordoba', 2, 1515411964, 20);
INSERT INTO "main"."provincia" VALUES (56, 'Granada', 2, 1515411977, 20);
INSERT INTO "main"."provincia" VALUES (57, 'Huelva', 2, 1515411987, 20);
INSERT INTO "main"."provincia" VALUES (58, 'Jaen', 2, 1515412005, 20);
INSERT INTO "main"."provincia" VALUES (59, 'Malaga', 2, 1515412020, 20);
INSERT INTO "main"."provincia" VALUES (60, 'Sevilla', 2, 1515412032, 20);
INSERT INTO "main"."provincia" VALUES (61, 'Huesca', 2, 1515412046, 21);
INSERT INTO "main"."provincia" VALUES (62, 'Teruel', 2, 1515412059, 21);
INSERT INTO "main"."provincia" VALUES (63, 'Zaragoza', 2, 1515412075, 21);
INSERT INTO "main"."provincia" VALUES (64, 'Asturias', 2, 1515412106, 22);
INSERT INTO "main"."provincia" VALUES (65, 'Palma de Mallorca', 2, 1515412118, 23);
INSERT INTO "main"."provincia" VALUES (66, 'Las Palmas', 2, 1515412135, 24);
INSERT INTO "main"."provincia" VALUES (67, 'Tenerife', 2, 1515412152, 24);
INSERT INTO "main"."provincia" VALUES (68, 'Cantabria', 2, 1515412172, 25);
INSERT INTO "main"."provincia" VALUES (69, 'Albacete', 2, 1515412195, 26);
INSERT INTO "main"."provincia" VALUES (70, 'Ciudad Real', 2, 1515412215, 26);
INSERT INTO "main"."provincia" VALUES (71, 'Cuenca', 2, 1515412344, 26);
INSERT INTO "main"."provincia" VALUES (72, 'Guadalajara', 2, 1515412376, 26);
INSERT INTO "main"."provincia" VALUES (73, 'Toledo', 2, 1515412395, 26);
INSERT INTO "main"."provincia" VALUES (74, 'Avila', 2, 1515412413, 27);
INSERT INTO "main"."provincia" VALUES (75, 'Burgos', 2, 1515412429, 27);
INSERT INTO "main"."provincia" VALUES (76, 'Leon', 2, 1515412448, 27);
INSERT INTO "main"."provincia" VALUES (77, 'Palencia', 2, 1515412469, 27);
INSERT INTO "main"."provincia" VALUES (78, 'Salamanca', 2, 1515412489, 27);
INSERT INTO "main"."provincia" VALUES (79, 'Segovia', 2, 1515412499, 27);
INSERT INTO "main"."provincia" VALUES (80, 'Soria', 2, 1515412511, 27);
INSERT INTO "main"."provincia" VALUES (81, 'Valladolid', 2, 1515412547, 27);
INSERT INTO "main"."provincia" VALUES (82, 'Barcelona', 2, 1515412565, 28);
INSERT INTO "main"."provincia" VALUES (83, 'Gerona', 2, 1515412578, 28);
INSERT INTO "main"."provincia" VALUES (84, 'Lerida', 2, 1515412589, 28);
INSERT INTO "main"."provincia" VALUES (85, 'Tarragona', 2, 1515412604, 28);
INSERT INTO "main"."provincia" VALUES (86, 'Alicante', 2, 1515412636, 29);
INSERT INTO "main"."provincia" VALUES (87, 'Castellon', 2, 1515412647, 29);
INSERT INTO "main"."provincia" VALUES (88, 'Valencia', 2, 1515412657, 29);
INSERT INTO "main"."provincia" VALUES (89, 'Badajoz', 2, 1515412676, 30);
INSERT INTO "main"."provincia" VALUES (90, 'Caceres', 2, 1515412687, 30);
INSERT INTO "main"."provincia" VALUES (91, 'La Coruña', 2, 1515412702, 31);
INSERT INTO "main"."provincia" VALUES (92, 'Lugo', 2, 1515412717, 31);
INSERT INTO "main"."provincia" VALUES (93, 'Orense', 2, 1515412728, 31);
INSERT INTO "main"."provincia" VALUES (94, 'Pontevedra', 2, 1515412745, 31);
INSERT INTO "main"."provincia" VALUES (95, 'La Rioja', 2, 1515412761, 32);
INSERT INTO "main"."provincia" VALUES (96, 'Madrid', 2, 1515412780, 33);
INSERT INTO "main"."provincia" VALUES (97, 'Murcia', 2, 1515412792, 34);
INSERT INTO "main"."provincia" VALUES (98, 'Navarra', 2, 1515412802, 35);
INSERT INTO "main"."provincia" VALUES (99, 'Alava', 2, 1515412812, 36);
INSERT INTO "main"."provincia" VALUES (100, 'Guipuzcoa', 2, 1515412824, 36);
INSERT INTO "main"."provincia" VALUES (101, 'Vizcaya', 2, 1515412833, 36);
INSERT INTO "main"."provincia" VALUES (102, 'Ceuta', 2, 1515412842, 37);
INSERT INTO "main"."provincia" VALUES (103, 'Melilla', 2, 1515412860, 38);
INSERT INTO "main"."provincia" VALUES (104, 'Almeria', 2, 1515413248, 39);
INSERT INTO "main"."provincia" VALUES (105, 'Malaga', 2, 1515413275, 39);
INSERT INTO "main"."provincia" VALUES (106, 'Cadiz', 2, 1515413305, 39);
INSERT INTO "main"."provincia" VALUES (107, 'Palma de Mallorca', 2, 1515413373, 42);
INSERT INTO "main"."provincia" VALUES (108, 'Las Palmas', 2, 1515413473, 43);
INSERT INTO "main"."provincia" VALUES (109, 'Tenerife', 2, 1515413485, 43);
INSERT INTO "main"."provincia" VALUES (110, 'Barcelona', 2, 1515413511, 47);
INSERT INTO "main"."provincia" VALUES (111, 'Valencia', 2, 1515413531, 48);
INSERT INTO "main"."provincia" VALUES (112, 'Cartagena', 2, 1515413548, 53);

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
-- Records of publi
-- ----------------------------

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
-- Records of region
-- ----------------------------
INSERT INTO "main"."region" VALUES (1, 'Andalucia', 2, 1515410157, 1);
INSERT INTO "main"."region" VALUES (2, 'Aragon', 2, 1515410234, 1);
INSERT INTO "main"."region" VALUES (3, 'Asturias', 2, 1515410243, 1);
INSERT INTO "main"."region" VALUES (4, 'Islas Baleares', 2, 1515410257, 1);
INSERT INTO "main"."region" VALUES (5, 'Canarias', 2, 1515410268, 1);
INSERT INTO "main"."region" VALUES (6, 'Cantabria', 2, 1515410285, 1);
INSERT INTO "main"."region" VALUES (7, 'Castilla-La Mancha', 2, 1515410303, 1);
INSERT INTO "main"."region" VALUES (8, 'Castilla y Leon', 2, 1515410320, 1);
INSERT INTO "main"."region" VALUES (9, 'Cataluña', 2, 1515410330, 1);
INSERT INTO "main"."region" VALUES (10, 'Valencia', 2, 1515410352, 1);
INSERT INTO "main"."region" VALUES (11, 'Extremadura', 2, 1515410364, 1);
INSERT INTO "main"."region" VALUES (12, 'Galicia', 2, 1515410371, 1);
INSERT INTO "main"."region" VALUES (13, 'La Rioja', 2, 1515410383, 1);
INSERT INTO "main"."region" VALUES (14, 'Madrid', 2, 1515410393, 1);
INSERT INTO "main"."region" VALUES (15, 'Murcia', 2, 1515410405, 1);
INSERT INTO "main"."region" VALUES (16, 'Navarra', 2, 1515410418, 1);
INSERT INTO "main"."region" VALUES (17, 'Pais Vasco', 2, 1515410427, 1);
INSERT INTO "main"."region" VALUES (18, 'Ceuta', 2, 1515410449, 1);
INSERT INTO "main"."region" VALUES (19, 'Melilla', 2, 1515410460, 1);
INSERT INTO "main"."region" VALUES (20, 'Andalucia', 2, 1515410483, 2);
INSERT INTO "main"."region" VALUES (21, 'Aragon', 2, 1515410493, 2);
INSERT INTO "main"."region" VALUES (22, 'Asturias', 2, 1515410506, 2);
INSERT INTO "main"."region" VALUES (23, 'Islas Baleares', 2, 1515410523, 2);
INSERT INTO "main"."region" VALUES (24, 'Canarias', 2, 1515410530, 2);
INSERT INTO "main"."region" VALUES (25, 'Cantabria', 2, 1515410561, 2);
INSERT INTO "main"."region" VALUES (26, 'Castilla-La Mancha', 2, 1515410583, 2);
INSERT INTO "main"."region" VALUES (27, 'Castilla y Leon', 2, 1515410597, 2);
INSERT INTO "main"."region" VALUES (28, 'Cataluña', 2, 1515410605, 2);
INSERT INTO "main"."region" VALUES (29, 'Valencia', 2, 1515410613, 2);
INSERT INTO "main"."region" VALUES (30, 'Extremadura', 2, 1515410624, 2);
INSERT INTO "main"."region" VALUES (31, 'Galicia', 2, 1515410632, 2);
INSERT INTO "main"."region" VALUES (32, 'La Rioja', 2, 1515410645, 2);
INSERT INTO "main"."region" VALUES (33, 'Madrid', 2, 1515410659, 2);
INSERT INTO "main"."region" VALUES (34, 'Murcia', 2, 1515410666, 2);
INSERT INTO "main"."region" VALUES (35, 'Navarra', 2, 1515410680, 2);
INSERT INTO "main"."region" VALUES (36, 'Pais Vasco', 2, 1515410688, 2);
INSERT INTO "main"."region" VALUES (37, 'Ceuta', 2, 1515410707, 2);
INSERT INTO "main"."region" VALUES (38, 'Melilla', 2, 1515410714, 2);
INSERT INTO "main"."region" VALUES (39, 'Andalucia', 2, 1515410734, 3);
INSERT INTO "main"."region" VALUES (42, 'Islas Baleares', 2, 1515410771, 3);
INSERT INTO "main"."region" VALUES (43, 'Canarias', 2, 1515410778, 3);
INSERT INTO "main"."region" VALUES (47, 'Cataluña', 2, 1515410831, 3);
INSERT INTO "main"."region" VALUES (48, 'Valencia', 2, 1515410844, 3);
INSERT INTO "main"."region" VALUES (53, 'Murcia', 2, 1515410892, 3);

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "main"."sqlite_sequence";
CREATE TABLE sqlite_sequence(name,seq);

-- ----------------------------
-- Records of sqlite_sequence
-- ----------------------------
INSERT INTO "main"."sqlite_sequence" VALUES ('almacenes', 3);
INSERT INTO "main"."sqlite_sequence" VALUES ('pais', 3);
INSERT INTO "main"."sqlite_sequence" VALUES ('region', 57);
INSERT INTO "main"."sqlite_sequence" VALUES ('provincia', 112);
INSERT INTO "main"."sqlite_sequence" VALUES ('usuarios', 3);
INSERT INTO "main"."sqlite_sequence" VALUES ('publi', 0);
INSERT INTO "main"."sqlite_sequence" VALUES ('entidades', 3);
INSERT INTO "main"."sqlite_sequence" VALUES ('tiendas', 253);
INSERT INTO "main"."sqlite_sequence" VALUES ('mensaje', 0);

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
"last_connect"  INTEGER,
CONSTRAINT "fk_user" FOREIGN KEY ("creador_id") REFERENCES "usuarios" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
CONSTRAINT "fk_prov" FOREIGN KEY ("provincia_id") REFERENCES "provincia" ("id") ON DELETE CASCADE ON UPDATE CASCADE
);

-- ----------------------------
-- Records of tiendas
-- ----------------------------
INSERT INTO "main"."tiendas" VALUES (1, 'CA0421', 2, 1515489557, 56, 'Plgno .lnd.Mercagranada - 18015 - Granada', 958295782, '', 1515489557);
INSERT INTO "main"."tiendas" VALUES (2, 'CA1004', 2, 1515489627, 54, 'Avd.Alcalde M. De la Pinta, 21 - 11011 - Cadiz', 956250500, '', 1515489627);
INSERT INTO "main"."tiendas" VALUES (3, 'CA1005', 2, 1515489936, 54, 'Ctra.Sanlucar-Chipiona, km3 - 11540 - San Lucar', 956362156, '', 1515489936);
INSERT INTO "main"."tiendas" VALUES (4, 'CA1006', 2, 1515490086, 54, 'Plgno.Ind. Fabricas. C/Fresadores - 11100 - San Fernando', 956881665, '', 1515490086);
INSERT INTO "main"."tiendas" VALUES (5, 'CA1007', 2, 1515490249, 54, 'Plgno.Ind. Pelagatos. Avd.Bahia Algeciras - 11130 - Chiclana', 956533351, '', 1515490249);
INSERT INTO "main"."tiendas" VALUES (6, 'CA1008', 2, 1515490390, 58, 'Plgno.Ind. Los Jarales, s/n - 23700 - Linares', 953694504, '', 1515490390);
INSERT INTO "main"."tiendas" VALUES (7, 'CA2002', 2, 1515490595, 96, 'Avd.John Lenon, 1. Pol. Industrial Los Angeles - 28906 - Getafe', 914957322, '', 1515490595);
INSERT INTO "main"."tiendas" VALUES (8, 'CA2003', 2, 1515490701, 97, 'Paraje Camino de la Estacion, s/n - 30400 - Caravaca', 968705630, '', 1515490701);
INSERT INTO "main"."tiendas" VALUES (9, 'CA2004', 2, 1515492089, 96, 'Ctra.Villaverde-Vallecas Km.3,8 Parc-3 Fase-2ª - 28053 - Mercamadrid', 915076906, '', 1515492089);
INSERT INTO "main"."tiendas" VALUES (10, 'SU0035', 2, 1515492912, 7, 'Polig. Industrial El Fuerte C/ Genal s/n - 29400 - Ronda', 952872643, '', 1515492912);
INSERT INTO "main"."tiendas" VALUES (11, 'SU1082', 2, 1515493031, 2, 'AV. Cruz Roja s/n (C.C. Merca 80, San Benito) - 11407 - Jerez', ' 956306750', '', 1515493031);
INSERT INTO "main"."tiendas" VALUES (12, 'SU1177', 2, 1515493603, 2, 'Ctra. Chipiona km 0,5 - 11540 - San Lucar', 956367878, '', 1515493603);
INSERT INTO "main"."tiendas" VALUES (13, 'SU2105', 2, 1515493739, 20, 'C/ Mejico, 33 Pol. Indu. El Balconcillo - 19004 - Guadalajara', 949208520, '', 1515493739);
INSERT INTO "main"."tiendas" VALUES (14, 'SU3074', 2, 1515493958, 8, 'Avda. Juan XXIII, s/n - 41710 - Utrera', 955861953, '', 1515493958);
INSERT INTO "main"."tiendas" VALUES (15, 'SU0001', 2, 1515494108, 7, 'Avda. Pío Baroja, 6 - 29017 - Málaga', 952296550, '', 1515494108);
INSERT INTO "main"."tiendas" VALUES (16, 'SU0004', 2, 1515494246, 7, 'Tomas Escalonilla, 4 (Gamarra) - 29010 - Málaga', 952284700, '', 1515494246);
INSERT INTO "main"."tiendas" VALUES (17, 'SU0005', 2, 1515494361, 7, 'Juan Antonio Tercero, 5 (Miraflores) - 29011 - Málaga', 952276262, '', 1515494361);
INSERT INTO "main"."tiendas" VALUES (18, 'SU0007', 2, 1515494479, 7, 'Reding, 10 - 29016 - Málaga', 952212859, '', 1515494479);
INSERT INTO "main"."tiendas" VALUES (19, 'SU0009', 2, 1515494616, 7, 'Crta.Cádiz Km.-196 - 29649 - Mijas-Costa', 952932961, '', 1515494616);
INSERT INTO "main"."tiendas" VALUES (20, 'SU0010', 2, 1515499681, 7, 'Avda.Principal del Candado, 2 - 29018 - Málaga', 952200387, '', 1515499681);
INSERT INTO "main"."tiendas" VALUES (21, 'SU0012', 2, 1515499855, 7, 'Ayala, 92 - 29002 - Málaga', 952212859, '', 1515499855);
INSERT INTO "main"."tiendas" VALUES (22, 'SU0014', 2, 1515499948, 7, 'Urb.La Capellania, Manzana-9 - 29130 - Alh. La Torre', 952417031, '', 1515499948);
INSERT INTO "main"."tiendas" VALUES (23, 'SU0016', 2, 1515500040, 7, 'Paseo del Colorado, 25 - 29620 - Torremolinos', 952382450, '', 1515500040);
INSERT INTO "main"."tiendas" VALUES (24, 'SU0017', 2, 1515500097, 7, 'Cervantes, 5 - 29016 - Málaga', 952603747, '', 1515500097);
INSERT INTO "main"."tiendas" VALUES (25, 'SU0018', 2, 1515500159, 7, 'Urb. Bell-Air, Camino Calderón - 29680 - Estepona', 952887171, '', 1515500159);
INSERT INTO "main"."tiendas" VALUES (26, 'SU0020', 2, 1515500251, 4, 'Avda.Peroné, Parcela-3 - 18680 - Salobreña', 958612234, '', 1515500251);
INSERT INTO "main"."tiendas" VALUES (27, 'SU0021', 2, 1515500331, 4, 'Pontanilla, 1 - 18680 - Salobreña', 958611258, '', 1515500331);
INSERT INTO "main"."tiendas" VALUES (28, 'SU0022', 2, 1515500380, 4, 'Rambla de Capuchinos, 8 - 18600 - Motril', 958605401, '', 1515500380);
INSERT INTO "main"."tiendas" VALUES (29, 'SU0023', 2, 1515500434, 4, 'Ancha, 63 - 18600 - Motril', 958611258, '', 1515500434);
INSERT INTO "main"."tiendas" VALUES (30, 'SU0024', 2, 1515500493, 4, 'Julio Moreno, 57 - 18600 - Motril', 958820335, '', 1515500493);
INSERT INTO "main"."tiendas" VALUES (31, 'SU0025', 2, 1515500554, 4, 'Avda.Salobreña, 23 - 18600 - Motril', 958600872, '', 1515500554);
INSERT INTO "main"."tiendas" VALUES (32, 'SU0026', 2, 1515500604, 4, 'Rodriguez Acosta, 11 - 18600 - Motril', 958604073, '', 1515500604);
INSERT INTO "main"."tiendas" VALUES (33, 'SU0031', 2, 1515500676, 7, 'Plaza Antonio Estrada, s/n - 29720 - La Cala del Moral', 952408044, '', 1515500676);
INSERT INTO "main"."tiendas" VALUES (34, 'SU0034', 2, 1515500762, 7, 'Avda. Santos Reing,5 - 29640 - Fuengirola', 952588860, '', 1515500762);
INSERT INTO "main"."tiendas" VALUES (35, 'SU0036', 2, 1515500856, 7, 'Antonio Ferrandis - 29780 - Nerja', 952527290, '', 1515500856);
INSERT INTO "main"."tiendas" VALUES (36, 'SU0038', 2, 1515500938, 7, 'C/ Gordon, 13 - 29013 - Málaga', 952254166, '', 1515500938);
INSERT INTO "main"."tiendas" VALUES (37, 'SU0040', 2, 1515501006, 7, 'Martinez Maldonado, 68 - 29010 - Málaga', 952396366, '', 1515501006);
INSERT INTO "main"."tiendas" VALUES (38, 'SU0044', 2, 1515548426, 7, 'Avda.de Carlos Haya, 60 - 29010 - Málaga', 952308662, '', 1515548426);
INSERT INTO "main"."tiendas" VALUES (39, 'SU0045', 2, 1515548484, 7, 'Camino de los Guindos, 8 - 29004 - Málaga', 952235397, '', 1515548484);
INSERT INTO "main"."tiendas" VALUES (40, 'SU0049', 2, 1515548620, 7, 'Cuarteles, 49 - 29002 - Málaga', 952353659, '', 1515548620);
INSERT INTO "main"."tiendas" VALUES (41, 'SU0050', 2, 1515548701, 7, 'Avda. de los Manantiales - 29620 - Torremolinos', 952389000, '', 1515548701);
INSERT INTO "main"."tiendas" VALUES (42, 'SU0052', 2, 1515548813, 7, 'Urb. Cortijo Blanco, Parcela A-1 - 29670 - S.Pedro de Alcantara', 952787500, '', 1515548813);
INSERT INTO "main"."tiendas" VALUES (43, 'SU0055', 2, 1515548928, 7, 'C/ Atarazanas, s/n - 29005 - Málaga', 952210725, '', 1515548928);
INSERT INTO "main"."tiendas" VALUES (44, 'SU0056', 2, 1515548983, 7, 'Avd/ Manolete s/n (CC. Nueva Andalucía, L.43) - 29660 - Marbella', 952813040, '', 1515548983);
INSERT INTO "main"."tiendas" VALUES (45, 'SU0057', 2, 1515549080, 7, 'Avda. Antonio Machado, 20 - 29639 - Benalmadena', 952560773, '', 1515549080);
INSERT INTO "main"."tiendas" VALUES (46, 'SU0058', 2, 1515549125, 7, 'Casablanca, 5 - 29620 - Torremolinos', 952376648, '', 1515549125);
INSERT INTO "main"."tiendas" VALUES (47, 'SU0061', 2, 1515549191, 7, 'Ctra. Benagalbon S/N - 29730 - Rincón de la Victoria', 952971622, '', 1515549191);
INSERT INTO "main"."tiendas" VALUES (48, 'SU0076', 2, 1515549252, 7, 'Ctra.Cadiz Km 192 Urb. Elviria - 29600 - Marbella', 952839051, '', 1515549252);
INSERT INTO "main"."tiendas" VALUES (49, 'SU0077', 2, 1515549322, 7, 'Urb. Guadalmina Alta C.C 4 fase - 29670 - Marbella', 952886165, '', 1515549322);
INSERT INTO "main"."tiendas" VALUES (50, 'SU0083', 2, 1515549413, 7, 'Avda. Salvador Vicente s/n - 29631 - Arroyo la miel', 952441340, '', 1515549413);
INSERT INTO "main"."tiendas" VALUES (51, 'SU0086', 2, 1515549490, 7, 'Urb. La Portada Parcela 4 D - 29680 - Estepona', 952801694, '', 1515549490);
INSERT INTO "main"."tiendas" VALUES (52, 'SU0103', 2, 1515549553, 7, 'Ctra. N. 340 Km. 179 Nagüeles - 29602 - Marbella', 952824206, '', 1515549553);
INSERT INTO "main"."tiendas" VALUES (53, 'SU0106', 2, 1515549595, 7, 'Avda.Juan Sebastian Elcano, 30 - 29017 - Málaga', 952298132, '', 1515549595);
INSERT INTO "main"."tiendas" VALUES (54, 'SU0115', 2, 1515549692, 7, 'Crta. Sabinillas Manilva Km 0,2 - 29691 - Minalva', 952891478, '', 1515549692);
INSERT INTO "main"."tiendas" VALUES (55, 'SU0129', 2, 1515549743, 7, 'Nuevo Boulevard, s/n - 29649 - Mijas/Costa', 952599058, '', 1515549743);
INSERT INTO "main"."tiendas" VALUES (56, 'SU0133', 2, 1515549808, 7, 'Almenara - 29680 - Estepona', 951316812, '', 1515549808);
INSERT INTO "main"."tiendas" VALUES (57, 'SU0310', 2, 1515549866, 7, 'Villafuerte, 35 - 29017 - Málaga', 952400410, '', 1515549866);
INSERT INTO "main"."tiendas" VALUES (58, 'SU0311', 2, 1515550187, 7, 'Paseo Cerrado de Calderón edif. Júpiter - 29018 - Málaga', 952200017, '', 1515550187);
INSERT INTO "main"."tiendas" VALUES (59, 'SU0314', 2, 1515550246, 7, 'Ctra. Almería, Finca La Batería - 29740 - Torre del Mar', 952542908, '', 1515550246);
INSERT INTO "main"."tiendas" VALUES (60, 'SU0315', 2, 1515550320, 7, 'Ctra. Almería, km. 283 - 29770 -Torrox', 952531149, '', 1515550320);
INSERT INTO "main"."tiendas" VALUES (61, 'SU0317', 2, 1515550450, 7, 'Avda. Mayorazgo - 29016 - Málaga', 952220634, '', 1515550450);
INSERT INTO "main"."tiendas" VALUES (62, 'SU0319', 2, 1515550505, 7, 'C/ Ronda s/n - 29730 - Rincón de la Victoria', 952402216, '', 1515550505);
INSERT INTO "main"."tiendas" VALUES (63, 'SU0322', 2, 1515550544, 7, 'C/ Pintor Joaquín Sorolla 7 - 29016 - Málaga', 952228343, '', 1515550544);
INSERT INTO "main"."tiendas" VALUES (64, 'SU0324', 2, 1515550590, 7, 'C/ Princesa, 6, Edif. Boulevard - 29740 - Torre del Mar', 952544112, '', 1515550590);
INSERT INTO "main"."tiendas" VALUES (65, 'SU0326', 2, 1515550638, 7, 'Urb. El Capistrano - 29780 - Nerja', 952520047, '', 1515550638);
INSERT INTO "main"."tiendas" VALUES (66, 'SU0327', 2, 1515550679, 7, 'C/ Dr. Fleming, 26 - 29740 - Torre del Mar', 952540441, '', 1515550679);
INSERT INTO "main"."tiendas" VALUES (67, 'SU0331', 2, 1515550732, 7, 'C/ Camino Viejo de Málaga - 29700 - Velez-Málaga', 952506500, '', 1515550732);
INSERT INTO "main"."tiendas" VALUES (68, 'SU0333', 2, 1515550775, 7, 'C/ Maestro Antonio Márquez - 29749 - Almayete', 952556463, '', 1515550775);
INSERT INTO "main"."tiendas" VALUES (69, 'SU0334', 2, 1515550870, 7, 'Paseo de los Tilos, 53 - 29006 - Málaga', 952313850, '', 1515550870);
INSERT INTO "main"."tiendas" VALUES (70, 'SU0336', 2, 1515550907, 7, 'Ctra. Almería, 19-21 - 29017 - Málaga', 952297101, '', 1515550907);
INSERT INTO "main"."tiendas" VALUES (71, 'SU0338', 2, 1515551008, 7, 'C/ Tamayo y Baus, 2 C. Sta. Inés - 29010 - Málaga', 952615578, '', 1515551008);
INSERT INTO "main"."tiendas" VALUES (72, 'SU0341', 2, 1515551057, 7, 'C/ Valentuñana, ed. Marbeland - 29601 - Marbella', 952770835, '', 1515551057);
INSERT INTO "main"."tiendas" VALUES (73, 'SU0342', 2, 1515551109, 7, 'Pza. Juan Macías - 29670 - S.Pedro de Alcántara', 952782484, '', 1515551109);
INSERT INTO "main"."tiendas" VALUES (74, 'SU0343', 2, 1515551549, 7, 'Avda. Gral. López Dominguez, 19 - 29603 - Marbella', 952778793, '', 1515551549);
INSERT INTO "main"."tiendas" VALUES (75, 'SU0344', 2, 1515551668, 7, 'Avda. Manilva, 1 - 29692 - San Luis de Sabinillas', 952890834, '', 1515551668);
INSERT INTO "main"."tiendas" VALUES (76, 'SU0348', 2, 1515551706, 7, 'Ctra. Cádiz km 196 - 29649 - Mijas-Costa', 952930160, '', 1515551706);
INSERT INTO "main"."tiendas" VALUES (77, 'SU0349', 2, 1515551741, 7, 'Ctra. Cádiz km 188 Urb. El Rosario - 29600 - Marbella', 952832036, '', 1515551741);
INSERT INTO "main"."tiendas" VALUES (78, 'SU0351', 2, 1515551819, 7, 'Avda. Obispo Herrera Oria - 29631 - Arroyo la Miel', 952564931, '', 1515551819);
INSERT INTO "main"."tiendas" VALUES (79, 'SU0352', 2, 1515551863, 7, 'Avda. Benalmádena, esq. A.Soler - 29620 - Torremolinos', 952053670, '', 1515551863);
INSERT INTO "main"."tiendas" VALUES (80, 'SU0355', 2, 1515551909, 7, 'CC Los Olivos - 29649 - Mijas-Costa', 952904870, '', 1515551909);
INSERT INTO "main"."tiendas" VALUES (81, 'SU0402', 2, 1515551952, 4, 'Alhamar, 39 - 18004 - Granada', 958520459, '', 1515551952);
INSERT INTO "main"."tiendas" VALUES (82, 'SU0403', 2, 1515551983, 4, 'Emperatriz Eugenia, 12 - 18002 - Granada', 958289785, '', 1515551983);
INSERT INTO "main"."tiendas" VALUES (83, 'SU0404', 2, 1515552031, 4, 'Martínez Campos, 23 - 18002 - Granada', 958520421, '', 1515552031);
INSERT INTO "main"."tiendas" VALUES (84, 'SU0407', 2, 1515552063, 4, 'Santiago Lozano - 18011 - Granada', 958158076, '', 1515552063);
INSERT INTO "main"."tiendas" VALUES (85, 'SU0409', 2, 1515552126, 4, '1º de Mayo, 16 - 18140 - La Zubia', 958891135, '', 1515552126);
INSERT INTO "main"."tiendas" VALUES (86, 'SU0418', 2, 1515552196, 4, 'Hnos. Machado, ed. Cruz de Mayo, s/n - 18140 - La Zubia', 958592378, '', 1515552196);
INSERT INTO "main"."tiendas" VALUES (87, 'SU0423', 2, 1515552237, 4, 'Plaza de los Campos - 18009 - Granada', 958227841, '', 1515552237);
INSERT INTO "main"."tiendas" VALUES (88, 'SU0425', 2, 1515552282, 4, 'Medina Olmo s/n - 18500 - Guadix', 958662738, '', 1515552282);
INSERT INTO "main"."tiendas" VALUES (89, 'SU0427', 2, 1515552419, 1, 'C/ Juan Aparicio - 04500 - Fiñana', 950352443, '', 1515552419);
INSERT INTO "main"."tiendas" VALUES (90, 'SU0428', 2, 1515553982, 4, 'La Hispanidad - 18320 - Santa Fe', 958440301, '', 1515553982);
INSERT INTO "main"."tiendas" VALUES (91, 'SU1001', 2, 1515554244, 2, 'GARCIA CARRERA, 24 - 11009 - Cadiz', 956288803, '', 1515554244);
INSERT INTO "main"."tiendas" VALUES (92, 'SU1002', 2, 1515554306, 2, 'ANTONIO MUÑOZ QUERO, 2 - 11012 - Cadiz', 956264653, '', 1515554306);
INSERT INTO "main"."tiendas" VALUES (93, 'SU1003', 2, 1515554332, 2, 'DIEGO ARIAS, 9-15 - 11002 - Cadiz', 956228876, '', 1515554332);
INSERT INTO "main"."tiendas" VALUES (94, 'SU1004', 2, 1515554379, 2, 'FACTORIA DE MATAGORDA,S/N - 11510 - Puerto Real', 956835062, '', 1515554379);
INSERT INTO "main"."tiendas" VALUES (95, 'SU1005', 2, 1515554413, 2, 'JUAN SEBASTIAN ELCANO,S/N - 11520 - Rota', 956815015, '', 1515554413);
INSERT INTO "main"."tiendas" VALUES (96, 'SU1006', 2, 1515554440, 2, 'SAN RAFAEL, 14 - 11520 - Rota', 956811159, '', 1515554440);
INSERT INTO "main"."tiendas" VALUES (97, 'SU1009', 2, 1515554472, 2, 'AVD. ANDALUCIA, 66 - 11007 - Cadiz', 956261703, '', 1515554472);
INSERT INTO "main"."tiendas" VALUES (98, 'SU1012', 2, 1515554510, 2, 'AVD. DE LA DIPUTACION - 11130 - Chiclana', 956532440, '', 1515554510);
INSERT INTO "main"."tiendas" VALUES (99, 'SU1014', 2, 1515554544, 2, 'MAGISTRAL CABRERA - 11130 - Chiclana', 956533888, '', 1515554544);
INSERT INTO "main"."tiendas" VALUES (100, 'SU1019', 2, 1515554635, 2, 'NUEVA, 28 - 11510 - Puerto Real', 956834455, '', 1515554635);
INSERT INTO "main"."tiendas" VALUES (101, 'SU1020', 2, 1515554666, 2, 'GARCIA GAMERO, 28 - 11012 - Cadiz', 956250814, '', 1515554666);
INSERT INTO "main"."tiendas" VALUES (102, 'SU1026', 2, 1515554701, 2, 'CRUZ DEL MONAGUILLO, 7 - 11540 - San Lucar', 956367150, '', 1515554701);
INSERT INTO "main"."tiendas" VALUES (103, 'SU1027', 2, 1515554741, 2, 'SANTO DOMINGO, 10 - 11402 - Jerez', 956338252, '', 1515554741);
INSERT INTO "main"."tiendas" VALUES (104, 'SU1028', 2, 1515554823, 2, 'San Francisco, 7 - 11203 - Algeciras', 956634275, '', 1515554823);
INSERT INTO "main"."tiendas" VALUES (105, 'SU1029', 2, 1515554854, 2, 'José Santacana, 5 - 11201 - Algeciras', 956632442, '', 1515554854);
INSERT INTO "main"."tiendas" VALUES (106, 'SU1032', 2, 1515554987, 2, 'PLAZA ALBORAN (Bda. Nazaret) - 11406 - Jerez', 956333111, '', 1515554987);
INSERT INTO "main"."tiendas" VALUES (107, 'SU1033', 2, 1515555033, 2, 'Avda. Bruselas - 11205 - Algeciras', 956630054, '', 1515555033);
INSERT INTO "main"."tiendas" VALUES (108, 'SU1037', 2, 1515555090, 2, 'AVD. DE LA LIBERTAD - 11500 - Pto.Sta.María', 956873365, '', 1515555090);
INSERT INTO "main"."tiendas" VALUES (109, 'SU1038', 2, 1515555134, 2, 'Avd/ Juan Carlos I s/n (Estadio Carranza) - 11010 - Cadiz', 956256569, '', 1515555134);
INSERT INTO "main"."tiendas" VALUES (110, 'SU1039', 2, 1515555167, 2, 'CC TARTESSUS, NOVO SANCTI PETRI - 11130 - Chiclana', 956494081, '', 1515555167);
INSERT INTO "main"."tiendas" VALUES (111, 'SU1065', 2, 1515555197, 2, 'AVD.ANDALUCIA, 9 - 11160 - Barbate', 956432101, '', 1515555197);
INSERT INTO "main"."tiendas" VALUES (112, 'SU1066', 2, 1515555837, 2, 'C/ ALONSO CANO, 1 - 11010 - Cadiz', 956273758, '', 1515555837);
INSERT INTO "main"."tiendas" VALUES (113, 'SU1070', 2, 1515555862, 2, 'C/ NARDO, 1 - 11520 - Rota', 956811402, '', 1515555862);
INSERT INTO "main"."tiendas" VALUES (114, 'SU1072', 2, 1515555901, 2, 'C/ CAÑUELO, 1 - 11190 - Benalup', 956424405, '', 1515555901);
INSERT INTO "main"."tiendas" VALUES (115, 'SU1079', 2, 1515556048, 2, 'Plaza GUERRA JIMENEZ S/N - 11001 - Cadiz', 956226300, '', 1515556048);
INSERT INTO "main"."tiendas" VALUES (116, 'SU1089', 2, 1515556085, 2, 'Ruiz Zorrilla, 21-22 - 11203 - Algeciras', 956651450, '', 1515556085);
INSERT INTO "main"."tiendas" VALUES (117, 'SU1090', 2, 1515556121, 2, 'AVD.INDEPENDENCIA,67 - 11580 - S.Jose Valle', 956160785, '', 1515556121);
INSERT INTO "main"."tiendas" VALUES (118, 'SU1091', 2, 1515556166, 2, 'C/ CUARTEL S/N - 11570 - Barca Florida', 956390300, '', 1515556166);
INSERT INTO "main"."tiendas" VALUES (119, 'SU1094', 2, 1515556198, 2, 'C/ Alfonso XI - 11201 - Algeciras', 956657957, '', 1515556198);
INSERT INTO "main"."tiendas" VALUES (120, 'SU1105', 2, 1515556231, 2, 'AVD.ANDALUCIA, Esq.Avd.Buenavista - 11150 - Vejer', 956451631, '', 1515556231);
INSERT INTO "main"."tiendas" VALUES (121, 'SU1107', 2, 1515556267, 2, 'Hnos.Lauhle,Esq.Daoiz,Esq.Lopez Rdguez - 11100 - San Fernando', 956593082, '', 1515556267);
INSERT INTO "main"."tiendas" VALUES (122, 'SU1108', 2, 1515556352, 2, 'Urb.El Aguila C/ Mochuelo, esq.Buho Real - 11500 - Pto.Sta.María', 956854209, '', 1515556352);
INSERT INTO "main"."tiendas" VALUES (123, 'SU1110', 2, 1515556384, 2, 'C/ San García - 11207 - Algeciras', 956603402, '', 1515556384);
INSERT INTO "main"."tiendas" VALUES (124, 'SU1114', 2, 1515556423, 2, 'URB. LOS GALLOS - 11130 - Chiclana', 956497401, '', 1515556423);
INSERT INTO "main"."tiendas" VALUES (125, 'SU1119', 2, 1515556730, 2, 'C/ San Pablo - 11300 - La Linea de la Concepción', 956177105, '', 1515556730);
INSERT INTO "main"."tiendas" VALUES (126, 'SU1126', 2, 1515556772, 2, 'P.IND.CASA DE LA REINA parc.17 - 11520 - Rota', 956810725, '', 1515556772);
INSERT INTO "main"."tiendas" VALUES (127, 'SU1127', 2, 1515556805, 2, 'C/REAL, 24 - 11100 - San Fernando', 956591155, '', 1515556805);
INSERT INTO "main"."tiendas" VALUES (128, 'SU1134', 2, 1515556834, 2, 'AVD. DEL PERU S/N - 11007 - Cadiz', 956261154, '', 1515556834);
INSERT INTO "main"."tiendas" VALUES (129, 'SU1135', 2, 1515556863, 2, 'C/ JOSE RAMOS BORRERO - 11100 - San Fernando', 956898942, '', 1515556863);
INSERT INTO "main"."tiendas" VALUES (130, 'SU1146', 2, 1515556954, 2, 'AVD. SEVILLA Nº 10-12 - 11550 - Chipiona', 956373186, '', 1515556954);
INSERT INTO "main"."tiendas" VALUES (131, 'SU1148', 2, 1515556994, 2, 'AVD. GENERALISIMO S/N - 11570 - Barbate', 956430460, '', 1515556994);
INSERT INTO "main"."tiendas" VALUES (132, 'SU1156', 2, 1515557055, 2, 'C/ LOS TOREROS Nº 2 - 11500 - Pto.Sta.María', 956850474, '', 1515557055);
INSERT INTO "main"."tiendas" VALUES (133, 'SU1159', 2, 1515557086, 2, 'HUERTA DEL ROSARIO S/N - 11130 - Chiclana', 956535663, '', 1515557086);
INSERT INTO "main"."tiendas" VALUES (134, 'SU1165', 2, 1515557120, 2, 'URB. EL ALMENDRAL BL. 10 Y 11 - 11407 - Jerez', 956310601, '', 1515557120);
INSERT INTO "main"."tiendas" VALUES (135, 'SU2001', 2, 1515576427, 45, 'C/Caceres, 23 28045 Madrid - 28045 - Madrid', 914685150, '', 1515576427);
INSERT INTO "main"."tiendas" VALUES (136, 'SU2002', 2, 1515576490, 45, 'C/Marques de Jura Real,17 - 28019 - Madrid', 914693797, '', 1515576490);
INSERT INTO "main"."tiendas" VALUES (137, 'SU2004', 2, 1515576535, 45, 'C/ Galileo, 87 - 28003 - Madrid', 915542666, '', 1515576535);
INSERT INTO "main"."tiendas" VALUES (138, 'SU2007', 2, 1515576603, 45, 'Ctra. Collado-Villalba (Villalba) - 28400 - Madrid', 918518161, '', 1515576603);
INSERT INTO "main"."tiendas" VALUES (139, 'SU2101', 2, 1515576643, 45, 'Avda. Brasil,18 - 28028 - Madrid', 915970968, '', 1515576643);
INSERT INTO "main"."tiendas" VALUES (140, 'SU2102', 2, 1515576718, 45, 'C/ Juan Ramón Jimenez,41 - 28036 - Madrid', 913502639, '', 1515576718);
INSERT INTO "main"."tiendas" VALUES (141, 'SU2104', 2, 1515576751, 45, 'C/ Arturo Soria, 310 - 28033 - Madrid', 917672710, '', 1515576751);
INSERT INTO "main"."tiendas" VALUES (142, 'SU2107', 2, 1515576811, 45, 'C/ Batalla Bailen, 2 - 28400 - Collado-Villaba (Madrid)', 918508234, '', 1515576811);
INSERT INTO "main"."tiendas" VALUES (143, 'SU2109', 2, 1515576872, 45, 'C/ Jesusa Lara y Fuente Albadalejo s/n - 28250 - Torrelodones', 918593402, '', 1515576872);
INSERT INTO "main"."tiendas" VALUES (144, 'SU2116', 2, 1515576919, 45, 'Avda. Dos Castillas, 16 - 28223 - Pozuelo', 913515256, '', 1515576919);
INSERT INTO "main"."tiendas" VALUES (145, 'SU2119', 2, 1515576951, 45, 'C/ Badajoz,13 - 28027 - Madrid', 913266179, '', 1515576951);
INSERT INTO "main"."tiendas" VALUES (146, 'SU2123', 2, 1515577015, 45, 'C/ Alcala, 281 - 28027 - Madrid', 914036081, '', 1515577015);
INSERT INTO "main"."tiendas" VALUES (147, 'SU2125', 2, 1515577403, 45, 'C/ Illescas, 185 - 28047 - Madrid', 914463864, '', 1515577403);
INSERT INTO "main"."tiendas" VALUES (148, 'SU2126', 2, 1515577488, 45, 'C/ Illescas, 197 - 28047  - Madrid', 917193966, '', 1515577488);
INSERT INTO "main"."tiendas" VALUES (149, 'SU2501', 2, 1515577641, 21, 'Travesía Marqués de Mendigorría, 4 - 45007 - Toledo', 925215622, '', 1515577641);
INSERT INTO "main"."tiendas" VALUES (150, 'SU2508', 2, 1515577739, 22, 'C/ De la Huerta, 4 - 05400 - Arenas S. Pedro', 920371390, '', 1515577739);
INSERT INTO "main"."tiendas" VALUES (151, 'SU2519', 2, 1515577802, 45, 'C/ Ortega y Gasset, 63 - 28006 - Madrid', 914023577, '', 1515577802);
INSERT INTO "main"."tiendas" VALUES (152, 'SU2521', 2, 1515577830, 45, 'C/ Golfo de Salonia, 2 local  - 28033 - Madrid', 913839479, '', 1515577830);
INSERT INTO "main"."tiendas" VALUES (153, 'SU2522', 2, 1515578439, 45, 'C/ Julián Romea, 4 - 28003 - Madrid', 915535930, '', 1515578439);
INSERT INTO "main"."tiendas" VALUES (154, 'SU2528', 2, 1515578496, 21, 'C/ Trinidad, 43-45 Talavera R - 45600 - Toledo', 925722553, '', 1515578496);
INSERT INTO "main"."tiendas" VALUES (155, 'SU2532', 2, 1515578611, 45, 'C/ Betanzos, 87 - 28034 - Madrid', 917385754, '', 1515578611);
INSERT INTO "main"."tiendas" VALUES (156, 'SU2533', 2, 1515578649, 45, 'C/ Fermin Caballero, 10 - 28019 - Madrid', 917392226, '', 1515578649);
INSERT INTO "main"."tiendas" VALUES (157, 'SU2536', 2, 1515578716, 45, 'C/ Martin Iriarte s/n Las Matas - 28230 - Las Rozas (Madrid)', 916304477, '', 1515578716);
INSERT INTO "main"."tiendas" VALUES (158, 'SU2537', 2, 1515578751, 45, 'C/ Cañada Nueva, 2 - 28200 - San Lorenzo del Escorial', 918901172, '', 1515578751);
INSERT INTO "main"."tiendas" VALUES (159, 'SU2539', 2, 1515578787, 45, 'C/ Sanchez Barcaiztegui, 41-43 - 28007 - Madrid', 914330251, '', 1515578787);
INSERT INTO "main"."tiendas" VALUES (160, 'SU2543', 2, 1515578820, 45, 'C/ Bolivar, 15 - 28045 - Madrid', 915279508, '', 1515578820);
INSERT INTO "main"."tiendas" VALUES (161, 'SU2545', 2, 1515582439, 45, 'C/ Santa Virgilia, 6 - 28033 - Madrid', 917630844, '', 1515582439);
INSERT INTO "main"."tiendas" VALUES (162, 'SU2546', 2, 1515582494, 45, 'C/ Del Rio s/n C.C. - 28972 - Puerta de Miraflores de la Sierra', 918444224, '', 1515582494);
INSERT INTO "main"."tiendas" VALUES (163, 'SU2547', 2, 1515582521, 45, 'C/ Cartagena, 20 - 28028 - Madrid', 913677899, '', 1515582521);
INSERT INTO "main"."tiendas" VALUES (164, 'SU2548', 2, 1515582559, 45, 'C/ La Hermita, 14-16 - 28411 - Moralzarzal', 918577290, '', 1515582559);
INSERT INTO "main"."tiendas" VALUES (165, 'SU2549', 2, 1515582591, 45, 'Plaza Maria Minguez, 7 - 28470 - Cercedilla', 918521536, '', 1515582591);
INSERT INTO "main"."tiendas" VALUES (166, 'SU2550', 2, 1515582620, 45, 'C/ Martin de los Heros, 33-35 - 28008 - Madrid', 915484054, '', 1515582620);
INSERT INTO "main"."tiendas" VALUES (167, 'SU2553', 2, 1515582649, 45, 'Avda de Niza, s/n - 28022 - Madrid', 917754559, '', 1515582649);
INSERT INTO "main"."tiendas" VALUES (168, 'SU3003', 2, 1515582690, 8, 'Calle Divino Redentor, 7-11 - 41005 - Sevilla', 954634342, '', 1515582690);
INSERT INTO "main"."tiendas" VALUES (169, 'SU3005', 2, 1515582724, 8, 'Calle Asunción, 17 - 41011 - Sevilla', 954273428, '', 1515582724);
INSERT INTO "main"."tiendas" VALUES (170, 'SU3008', 2, 1515582752, 8, 'Calle Venecia, 19 - 41008 - Sevilla', 954414408, '', 1515582752);
INSERT INTO "main"."tiendas" VALUES (171, 'SU3011', 2, 1515582809, 8, 'Avda. Reina Mercedes, 45 - 41012 - Sevilla', 954238232, '', 1515582809);
INSERT INTO "main"."tiendas" VALUES (172, 'SU3012', 2, 1515583017, 8, 'Calle Guadalimar, s/n (Pedro Salvador) - 41013 - Sevilla', 954616098, '', 1515583017);
INSERT INTO "main"."tiendas" VALUES (173, 'SU3013', 2, 1515583050, 8, 'Avda. del Greco, 14 - 41007 - Sevilla', 954576694, '', 1515583050);
INSERT INTO "main"."tiendas" VALUES (174, 'SU3016', 2, 1515583086, 8, 'Calle Forjadores, 8 - 41015 - Sevilla', 954940976, '', 1515583086);
INSERT INTO "main"."tiendas" VALUES (175, 'SU3020', 2, 1515583132, 8, 'Calle Ingeniero la Cierva s/n - 41006 - Sevilla', 954650624, '', 1515583132);
INSERT INTO "main"."tiendas" VALUES (176, 'SU3032', 2, 1515583329, 8, 'Calle José Laguillo, 12 - 41003 - Sevilla', 954530948, '', 1515583329);
INSERT INTO "main"."tiendas" VALUES (177, 'SU3044', 2, 1515583370, 8, 'Avda. de Andalucía, s/n - 41309 - La Rinconada', 955791273, '', 1515583370);
INSERT INTO "main"."tiendas" VALUES (178, 'SU3055', 2, 1515583423, 8, 'Calle Marqués de Mina, 8 - 41002 - Sevilla', 954381011, '', 1515583423);
INSERT INTO "main"."tiendas" VALUES (179, 'SU3060', 2, 1515583456, 8, 'Plaza de la Alfalfa, 4-5 - 41004 - Sevilla', 954560092, '', 1515583456);
INSERT INTO "main"."tiendas" VALUES (180, 'SU3069', 2, 1515583505, 8, 'Calle Salado, 1-3 - 41010 - Sevilla', 954459638, '', 1515583505);
INSERT INTO "main"."tiendas" VALUES (181, 'SU3070', 2, 1515583538, 8, 'Calle Juan Nuñez, 33 - 41008 - Sevilla', 954435152, '', 1515583538);
INSERT INTO "main"."tiendas" VALUES (182, 'SU3093', 2, 1515583594, 8, 'Calle Cataluña, 3 - 41015 - Sevilla', 954373004, '', 1515583594);
INSERT INTO "main"."tiendas" VALUES (183, 'SU3160', 2, 1515583631, 5, 'C/ Alanis de la Sierra,11 - 21007 - Huelva', 959236745, '', 1515583631);
INSERT INTO "main"."tiendas" VALUES (184, 'SU3186', 2, 1515583671, 8, 'Calle Candelaria, 26 - 41006 - Sevilla', 954637139, '', 1515583671);
INSERT INTO "main"."tiendas" VALUES (185, 'SU3190', 2, 1515583711, 5, 'Avda. Galaroza - 21006 - Huelva', 959220492, '', 1515583711);
INSERT INTO "main"."tiendas" VALUES (186, 'SU3193', 2, 1515583759, 8, 'Avda. Rey de España (Ant. Crta. Sevilla) - 41807 - Espartina', 955711315, '', 1515583759);
INSERT INTO "main"."tiendas" VALUES (187, 'SU1315', 2, 1515583796, 2, 'CONIL', '-', '', 1515583796);
INSERT INTO "main"."tiendas" VALUES (188, 'SU1312', 2, 1515583830, 2, 'PLAYA DEL ROMPIDILLO S/N - Jerez', '-', '', 1515583830);
INSERT INTO "main"."tiendas" VALUES (189, 'SU1313', 2, 1515583944, 2, 'TARIFA', '-', '', 1515583944);
INSERT INTO "main"."tiendas" VALUES (190, 'SU1314', 2, 1515584694, 5, 'ALMONTE', '-', '', 1515584694);
INSERT INTO "main"."tiendas" VALUES (191, 'SU1316', 2, 1515584716, 5, 'ISLA CRISTINA', '-', '', 1515584716);
INSERT INTO "main"."tiendas" VALUES (192, 'SU1317', 2, 1515584768, 2, 'CONIL', '-', '', 1515584768);
INSERT INTO "main"."tiendas" VALUES (193, 'SU1311', 2, 1515584787, 5, 'AYAMONTE', '-', '', 1515584787);
INSERT INTO "main"."tiendas" VALUES (194, 'SU0054', 2, 1515584830, 52, 'C/ General Polavieja, s/n - 52006 - Melilla', 952678700, '', 1515584830);
INSERT INTO "main"."tiendas" VALUES (195, 'SU0092', 2, 1515584874, 52, 'C/ Madrid Edf. Zurbaran - 52005 - Melilla', 952675034, '', 1515584874);
INSERT INTO "main"."tiendas" VALUES (196, 'SU0124', 2, 1515584913, 52, 'C/Rafael Genil. Edificio San Lorenzo - 52001 - Melilla', 952676367, '', 1515584913);
INSERT INTO "main"."tiendas" VALUES (197, 'FR0001', 2, 1515584958, 51, 'Muelle Cañonero Dato - 51001 - Ceuta', 956507529, '', 1515584958);
INSERT INTO "main"."tiendas" VALUES (198, 'SU1202', 2, 1515584994, 51, 'Paseo del Revellín, 13 - 51001 - Ceuta', 956512910, '', 1515584994);
INSERT INTO "main"."tiendas" VALUES (199, 'SU1206', 2, 1515585026, 51, 'Paseo de Colón, 11 - 51001 - Ceuta', 956513823, '', 1515585026);
INSERT INTO "main"."tiendas" VALUES (200, 'SU1309', 2, 1515585068, 8, 'Avenida Felipe II, 12-14 - 41013 - Sevilla', 954238197, '', 1515585068);
INSERT INTO "main"."tiendas" VALUES (201, 'SU1310', 2, 1515585101, 8, 'Carretera de Carmona, 40 - 41008 - Sevilla', 954421759, '', 1515585101);
INSERT INTO "main"."tiendas" VALUES (202, 'SU2541', 2, 1515585166, 39, 'Avda. de las Angustias, s/n - 10300 - Navalmoral de la Mata', 927537206, '', 1515585166);
INSERT INTO "main"."tiendas" VALUES (203, 'SU1121', 2, 1515585313, 2, 'C/ Jose Carlos de Luna - 11202 - Algeciras', 956665654, '', 1515585313);
INSERT INTO "main"."tiendas" VALUES (204, 'SU2557', 2, 1515585378, 45, 'Calle Rafael Calvo 18 - 28010 - Madrid', 620614333, '', 1515585378);
INSERT INTO "main"."tiendas" VALUES (205, 'SU2558', 2, 1515585423, 45, 'Infanta María Teresa 12 - 28016 - Madrid', 91, '', 1515585423);
INSERT INTO "main"."tiendas" VALUES (206, 'SU8888', 2, 1515585468, 7, 'Jerez 16 - 29631 - Arroyo de la Miel', 696900910, '', 1515585468);
INSERT INTO "main"."tiendas" VALUES (207, 'SU2559', 2, 1515585567, 45, 'Bulevar Salvador Allende, nº 6 - 28108 - Alcobendas', 648053414, '', 1515585567);
INSERT INTO "main"."tiendas" VALUES (208, 'SU1318', 2, 1515585616, 2, 'Avd Alameda de Solano esq Avd de la libertad sn - 11130 - Chiclana', 648059581, '', 1515585616);
INSERT INTO "main"."tiendas" VALUES (209, 'SU2562', 2, 1515585662, 45, 'Calle - 28220 - Majadahonda', 91, '', 1515585662);
INSERT INTO "main"."tiendas" VALUES (210, 'SU1319', 2, 1515585761, 8, 'Calle Eduardo Dato 22 - 41018 - Sevilla', 954000000, '', 1515585761);
INSERT INTO "main"."tiendas" VALUES (211, 'SU2561', 2, 1515585820, 45, 'Avenida Ciudad de Barcelona 142 - 28007 - Madrid', 91000000, '', 1515585820);
INSERT INTO "main"."tiendas" VALUES (212, 'SU2564', 2, 1515585878, 45, 'Calle Alcala, 80 - 28009 - Madrid', 9100, '', 1515585878);
INSERT INTO "main"."tiendas" VALUES (213, 'CA9107', 2, 1515586465, 66, 'C/ CTRA AEROPUERTO KM 3 ZONA IND - Playa Honda', 928821194, '', 1515586465);
INSERT INTO "main"."tiendas" VALUES (214, 'CA9282', 2, 1515586516, 66, 'POLIGONO INDST.MONTAÑA ROJA NAVE 1,2,5, TERMINO N - Playa Blanca', 928519325, '', 1515586516);
INSERT INTO "main"."tiendas" VALUES (215, 'CA9105', 2, 1515586599, 66, 'C/ ARRECIFE, S/N POL. IND LOMO BLANCO - 35010- LAS TORRES', '928480832,27', '', 1515586599);
INSERT INTO "main"."tiendas" VALUES (216, 'CA9087', 2, 1515586691, 66, 'C/ PIZARRO 43 Y C/ PEJIN 17-POLIG. 11,09 - Corralejo', 928537373, '', 1515586691);
INSERT INTO "main"."tiendas" VALUES (217, 'CA9110', 2, 1515586735, 66, 'C/ PIZARRO, S/N EL CHARCO - Puerto Rosario', 928531353, '', 1515586735);
INSERT INTO "main"."tiendas" VALUES (218, 'CA9112', 2, 1515586780, 66, 'ZONA IND. PUERTO DE MORRO JABLE - Jandia', 928540465, '', 1515586780);
INSERT INTO "main"."tiendas" VALUES (219, 'CA9069', 2, 1515586852, 67, 'C/ LOMO BLANCO 42 STA BARBARA - Icod de los Vinos', 922813862, '', 1515586852);
INSERT INTO "main"."tiendas" VALUES (220, 'CA9070', 2, 1515586889, 67, 'POLIGONO IND. LAS ANDORRIÑAS C/ TAFETAN nave36 - San Miguel de Abona', 922735101, '', 1515586889);
INSERT INTO "main"."tiendas" VALUES (221, 'CA9072', 2, 1515633459, 67, 'PASEO LAS ARAUCARIAS, 40 - La Orotava', 922320248, '', 1515633459);
INSERT INTO "main"."tiendas" VALUES (222, 'CA9073', 2, 1515633507, 67, 'POLIGONO IND. SAN JERONIMO S/N - Las Arenas', 922320795, '', 1515633507);
INSERT INTO "main"."tiendas" VALUES (223, 'CA9170', 2, 1515633574, 67, 'P.I.. Bco. Las Torres SECTOR 10 Parc.8 Manz S-4 - Adeje', 922781236, '', 1515633574);
INSERT INTO "main"."tiendas" VALUES (224, 'CA9142', 2, 1515633638, 67, 'C/ EUROPA S/N-POL. AGROINDUSTRIAL BUENA VISTA - Breña Alta', 922429537, '', 1515633638);
INSERT INTO "main"."tiendas" VALUES (225, 'CA1009', 2, 1515633683, 103, 'Plgno.Ind. El Sepes. C/Dalia, S/N - 52006 - Melilla', 952675744, '', 1515633683);
INSERT INTO "main"."tiendas" VALUES (226, 'CA1290', 2, 1515633732, 102, 'Av.Muelle de Poniente, s/n - 51001 - Ceuta', 956501560, '', 1515633732);
INSERT INTO "main"."tiendas" VALUES (227, 'CA1300', 2, 1515634079, 102, ' Pg. El Tarajal nave 3 y 4 - 51101 - Ceuta', 956521305, '', 1515634079);
INSERT INTO "main"."tiendas" VALUES (228, 'CA9053', 2, 1515634157, 66, 'C/ GARCIA ESCAMEZ, 230 - TRASERA REPSOL - ARGANA A - 35500 - ARRECIFE_LPA', 928802028, '', 1515634157);
INSERT INTO "main"."tiendas" VALUES (229, 'CA9071', 2, 1515634324, 67, 'Poligono Industrial el Mayorazgo (Subida el mayor) - 38110 - El Mayorazgo', 922200179, '', 1515634324);
INSERT INTO "main"."tiendas" VALUES (230, 'CA9286', 2, 1515634393, 66, 'CALLE EL MODEM N 2224 POLIGONO INDUSTRIAL JINAMAR - 35220 - TELDE', 928710620, '', 1515634393);
INSERT INTO "main"."tiendas" VALUES (231, 'CA9109', 2, 1515634455, 66, 'CALLE ALCALDE ENRIQUE JORGE SN P IND - 35100 - PLAYA DEL INGLES', 619001105, '', 1515634455);
INSERT INTO "main"."tiendas" VALUES (232, 'CA9111', 2, 1515634544, 66, 'LOMO GUILLEN 21 - 35450 - SANTA MARIA DE GUIA', 669460880, '', 1515634544);
INSERT INTO "main"."tiendas" VALUES (233, 'CA9130', 2, 1515634587, 66, 'Calle Eucalipto parcela 6 - Arinaga - 35110 - Vecindario', 928793684, '', 1515634587);
INSERT INTO "main"."tiendas" VALUES (234, 'CA9102', 2, 1515634740, 66, 'Cuesta Ramon sn Merca Las Palmas - 35229 - Telde', 928711704, '', 1515634740);
INSERT INTO "main"."tiendas" VALUES (235, 'CA9104', 2, 1515634781, 66, 'MOLINO 2 - 35140 - ARGUINEGUIN', 928735076, '', 1515634781);
INSERT INTO "main"."tiendas" VALUES (236, 'CA9108', 2, 1515634821, 66, 'JUAN DOMINGUEZ PEREZ 38 PI EL SEBADAL - 35008 - GRAN CANARIA', 669460880, '', 1515634821);
INSERT INTO "main"."tiendas" VALUES (237, 'CA9103', 2, 1515634894, 66, 'SANTIAGO ASCANIO MONTEMAYOR 2 - 35008 - TELDE', 928685382, '', 1515634894);
INSERT INTO "main"."tiendas" VALUES (238, 'CA9101', 2, 1515634934, 66, 'AGUSTIN ITURBIDE 1 - 35200 - TELDE', 928706756, '', 1515634934);
INSERT INTO "main"."tiendas" VALUES (239, 'ACC JUAN J SISTER', 2, 1515636513, 107, 'Puerto de Mallorca - 07001 - Mallorca', 954000000, '', 1515636513);
INSERT INTO "main"."tiendas" VALUES (240, 'ACC LAS PALMAS GC', 2, 1515636579, 104, 'Puerto de Almeria - 04001 - Almeria', 954000000, '', 1515636579);
INSERT INTO "main"."tiendas" VALUES (241, 'ACC SOROLLA', 2, 1515636656, 105, 'Puerto de Malaga - 29003 - Malaga', 954000000, '', 1515636656);
INSERT INTO "main"."tiendas" VALUES (242, 'ACC CIUDAD MALAGA', 2, 1515636703, 106, 'Puerto de Algeciras - 11202 - Algeciras', 954000000, '', 1515636703);
INSERT INTO "main"."tiendas" VALUES (243, 'ACC MILENIUM 2', 2, 1515636750, 106, 'Puerto de Algeciras - 11202 - Algeciras', 677509283, '', 1515636750);
INSERT INTO "main"."tiendas" VALUES (244, 'ACC WISTERIA', 2, 1515636811, 104, 'Puerto de Almeria - 04001 - Almeria', 954000000, '', 1515636811);
INSERT INTO "main"."tiendas" VALUES (245, 'ACC ISABELLA', 2, 1515636840, 104, 'Puerto de Almeria - 04001 - Almeria', 954000000, '', 1515636840);
INSERT INTO "main"."tiendas" VALUES (246, 'ACC ALBAYZIN', 2, 1515636878, 106, 'Puerto de Cadiz - 11001 - Cadiz', 964000000, '', 1515636878);
INSERT INTO "main"."tiendas" VALUES (247, 'ACC ZURBARAN', 2, 1515636919, 107, 'Puerto de Palma - 07001 - Palma de Mallorca', 964000000, '', 1515636919);
INSERT INTO "main"."tiendas" VALUES (248, 'ACC TENACIA', 2, 1515637013, 110, 'Puerto de Barcelona - 08001 - Barcelona', 934000000, '', 1515637013);
INSERT INTO "main"."tiendas" VALUES (249, 'ACC FORTUNY', 2, 1515637044, 105, 'Puerto de Malaga - 29001 - Malaga', 934000000, '', 1515637044);
INSERT INTO "main"."tiendas" VALUES (250, 'ACC ALBORAN', 2, 1515637099, 106, 'Puerto de Algeciras - 11205 - Algeciras', 954000000, '', 1515637099);
INSERT INTO "main"."tiendas" VALUES (251, 'ACC DIMONIOS', 2, 1515637129, 107, 'Puerto de Palma - 07001 - Mallorca', 964000000, '', 1515637129);
INSERT INTO "main"."tiendas" VALUES (252, 'ACC FORZA', 2, 1515637174, 111, 'Puerto Valencia - 46021 - Valencia', 670462502, '', 1515637174);
INSERT INTO "main"."tiendas" VALUES (253, 'ACC ADRIATICO', 2, 1515637213, 110, 'Puerto Barcelona - 08027 - Barcelona', 617427767, '', 1515637213);

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

-- ----------------------------
-- Records of usuarios
-- ----------------------------
INSERT INTO "main"."usuarios" VALUES (1, 'admin', 'admin', 'admin', 'superusuario', 0, 0, '3F');
INSERT INTO "main"."usuarios" VALUES (2, 'Juan', 'Juan', 'jr2017', 'Juan Ramon', 1, 1, '1f');
INSERT INTO "main"."usuarios" VALUES (3, 'Matias', 'Matias', 1234, 'Matias Prats', 2, 2, '1f');
