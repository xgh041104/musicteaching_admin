/*
 Navicat Premium Dump SQL

 Source Server         : mysqllink
 Source Server Type    : MySQL
 Source Server Version : 80040 (8.0.40)
 Source Host           : localhost:3306
 Source Schema         : music_teaching

 Target Server Type    : MySQL
 Target Server Version : 80040 (8.0.40)
 File Encoding         : 65001

 Date: 07/07/2025 10:06:03
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for book
-- ----------------------------
DROP TABLE IF EXISTS `book`;
CREATE TABLE `book`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '书本Id，主键',
  `book_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '书本名称',
  `course_count` int NOT NULL DEFAULT 0 COMMENT '视频数',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `delete_at` datetime NULL DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 58 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '书籍信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of book
-- ----------------------------
INSERT INTO `book` VALUES (21, '七年级上册', 16, '2025-07-04 17:23:44', '2025-07-06 23:07:16', NULL);
INSERT INTO `book` VALUES (22, '七年级下册', 1, '2025-07-04 17:23:44', '2025-07-06 09:59:40', NULL);
INSERT INTO `book` VALUES (23, '八年级上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (24, '八年级下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (25, '九年级上册', 1, '2025-07-04 17:23:44', '2025-07-05 18:08:24', NULL);
INSERT INTO `book` VALUES (26, '九年级下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (27, '一年级 上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (28, '一年级 下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (29, '二年级 上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (30, '二年级 下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (31, '三年级 上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (32, '三年级 下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (33, '四年级 上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (34, '四年级 下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (35, '五年级 上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (36, '五年级 下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (37, '六年级 上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (38, '六年级 下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (39, '七年级上册', 9, '2025-07-04 17:23:44', '2025-07-06 21:17:36', NULL);
INSERT INTO `book` VALUES (40, '七年级下册', 7, '2025-07-04 17:23:44', '2025-07-06 17:39:36', NULL);
INSERT INTO `book` VALUES (41, '八年级上册', 4, '2025-07-04 17:23:44', '2025-07-06 17:24:00', NULL);
INSERT INTO `book` VALUES (42, '八年级下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (43, '九年级上册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (44, '九年级下册', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (45, '一年级上册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (46, '一年级下册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (47, '二年级上册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (48, '二年级下册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (49, '三年级上册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (50, '三年级下册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (51, '四年级上册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (52, '四年级下册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (53, '五年级上册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (54, '五年级下册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (55, '六年级上册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (56, '六年级下册（简谱）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);
INSERT INTO `book` VALUES (57, '七年级上册（新）', 0, '2025-07-04 17:23:44', '2025-07-04 17:23:44', NULL);

-- ----------------------------
-- Table structure for course
-- ----------------------------
DROP TABLE IF EXISTS `course`;
CREATE TABLE `course`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '课程id',
  `book_id` bigint UNSIGNED NOT NULL COMMENT '所属书本ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '课程标题',
  `video_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '视频路径',
  `record_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '录音路径',
  `summary` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '课程总结',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `delete_at` datetime NULL DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_book_id`(`book_id` ASC) USING BTREE,
  INDEX `idx_delete_at`(`delete_at` ASC) USING BTREE,
  CONSTRAINT `fk_course_book` FOREIGN KEY (`book_id`) REFERENCES `book` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 86 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of course
-- ----------------------------
INSERT INTO `course` VALUES (11, 21, '音乐课', 'static\\media\\2025\\07\\05\\ec5745786e3feb2c5ea24ba06a583b9b.mp4', 'static\\media\\2025\\07\\05\\ec5745786e3feb2c5ea24ba06a583b9b.mp3', '', '2025-07-05 15:10:16', '2025-07-05 15:10:16', NULL);
INSERT INTO `course` VALUES (12, 21, '音乐课', 'static\\media\\2025\\07\\05\\a9f9548a36ddcb6184ca36e817e93afc.mp4', 'static\\media\\2025\\07\\05\\a9f9548a36ddcb6184ca36e817e93afc.mp3', '', '2025-07-05 16:25:57', '2025-07-05 16:25:57', NULL);
INSERT INTO `course` VALUES (13, 21, '音乐课', 'static\\media\\2025\\07\\05\\b55070fd25f428566e377500139e3b7b.mp4', 'static\\media\\2025\\07\\05\\b55070fd25f428566e377500139e3b7b.mp3', '以下是针对音乐课的复习总结框架，供学生参考：\n\n---\n\n### **音乐课核心复习要点**\n1. **课程主题**  \n   - 基础乐理知识（如音符、节拍、音阶）  \n   - 乐器认知与演奏方法（如钢琴、吉他、打击乐器等）  \n   - 音乐欣赏（经典作品、作曲家风格、音乐流派）  \n   - 实践环节（合唱、节奏练习、简单作曲）  \n\n2. **重点知识回顾**  \n   - **乐理基础**：五线谱/简谱识谱、音符时值（全音符、二分音符等）、常用音乐术语（如“渐强”“快板”）。  \n   - **乐器学习**：不同乐器的发声原理、持琴/演奏姿势、基础指法（如吉他和弦、钢琴手型）。  \n   - **音乐史**：古典音乐（贝多芬、莫扎特）、民族音乐、现代流行音乐特点对比。  \n\n3. **课堂活动总结**  \n   - 分组演奏/合唱的协作经验  \n   - 节奏感训练（如拍手、打击乐器配合）  \n   - 音乐创作实践（如改编旋律、编写简单歌词）  \n\n4. **复习建议**  \n   - 每日练习乐器10分钟，巩固指法/技巧。  \n   - 听辨经典曲目片段，尝试分析其节奏、旋律特点。  \n   - 完成乐理小测试（如写出指定音符的时值、识别音阶）。  \n\n---\n\n### **温馨提示**  \n若需具体知识点梳理（如某类乐器操作、乐理难点），可提供详细内容，我将协助整理！ 😊  \n\n（注：若原提问中的“音乐课”内容有其他补充，请进一步说明，以便精准总结。）', '2025-07-05 16:45:00', '2025-07-05 16:45:00', NULL);
INSERT INTO `course` VALUES (14, 21, '音乐课', 'static\\media\\2025\\07\\05\\8dd12e1aa147a4ce5c426d61b11ab8ed.mp4', 'static\\media\\2025\\07\\05\\8dd12e1aa147a4ce5c426d61b11ab8ed.mp3', '1. **课程主题**  \n音乐课堂问候与人工智能认知  \n\n2. **主要知识点**  \n- 音乐课堂问候语的使用  \n- 人工智能的基本识别与提问  \n\n3. **难点提醒**  \n- 正确使用音乐课堂特定问候语  \n- 理解人工智能的基本概念及互动方式  \n\n4. **书面总结**  \n本节课学习了音乐课堂问候语及人工智能的基本认知，重点掌握问候礼仪与简单问答。', '2025-07-05 17:11:57', '2025-07-05 17:11:57', NULL);
INSERT INTO `course` VALUES (16, 21, '音乐课', 'static\\media\\2025\\07\\05\\c50e078bd4b3d21b6e38c57c92bdc6e9.mp4', 'static\\media\\2025\\07\\05\\c50e078bd4b3d21b6e38c57c92bdc6e9.mp3', '1. 课程主题  \n大学生防诈骗安全教育  \n\n2. 主要知识点  \n- 大学生常见诈骗类型：刷单诈骗、网贷诈骗、冒充熟人诈骗  \n- 诈骗高发原因：社会经验不足、轻信陌生人、急于求成心理  \n- 防骗关键措施：不点击陌生链接、不扫不明二维码、涉及钱财先核实身份、与师长协商  \n\n3. 难点提醒  \n- 混淆不同诈骗手段的识别特征  \n- 未充分核实对方身份即盲目转账  \n- 因急于获利忽视风险提示  \n\n4. 书面总结格式  \n【课程主题】大学生防诈骗安全教育  \n【核心内容】  \n① 数据警示：大学生被骗案件频发，刷单、网贷、冒充类诈骗占比超70%；  \n② 典型案例：虚假刷单诱导大额投入、冒充亲友紧急借款；  \n③ 受骗根源：社会经验匮乏、轻信心理、急功近利；  \n④ 防御要点：拒陌生链接/二维码、钱财事务核实身份、及时求助师长。  \n\n5. 总结内容  \n课程通过数据与案例分析大学生受骗现象，归纳诈骗类型及成因，强调防范需警惕陌生信息、核实身份、理性判断，避免因轻信或贪利陷入骗局。', '2025-07-05 17:18:19', '2025-07-05 17:18:19', NULL);
INSERT INTO `course` VALUES (17, 39, '11111', 'static\\media\\2025\\07\\05\\a1a6af35f31a404ac21eb99acee080ed.mp4', 'static\\media\\2025\\07\\05\\a1a6af35f31a404ac21eb99acee080ed.mp3', '1. 课程主题  \n2. 主要知识点  \n3. 难点提醒', '2025-07-05 17:21:35', '2025-07-05 17:21:35', NULL);
INSERT INTO `course` VALUES (18, 25, '1111', 'static\\media\\2025\\07\\05\\97d718f0a8917db1372f2a3775a0d644.mp4', 'static\\media\\2025\\07\\05\\97d718f0a8917db1372f2a3775a0d644.mp3', '1. 课程主题  \n数字与声音的对应关系认知  \n\n2. 主要知识点  \n- 数字“1”的书写与发音  \n- 拟声词“喂”的模拟与表达  \n- 简单重复符号的含义区分  \n\n3. 难点提醒  \n- 数字“1”与多个“1”组合的辨识  \n- “喂”的发音与实际呼叫场景的关联  \n- 符号“1111”与声音“喂喂喂”的对应逻辑', '2025-07-05 18:08:24', '2025-07-05 18:08:24', NULL);
INSERT INTO `course` VALUES (19, 41, '45456465', 'static\\media\\2025\\07\\05\\367a0f2339f406a76e6e11661b136826.mp4', 'static\\media\\2025\\07\\05\\367a0f2339f406a76e6e11661b136826.mp3', '1. 课程主题  \n未明确提及  \n\n2. 主要知识点  \n未明确提及  \n\n3. 难点提醒  \n未明确提及  \n\n4. 比较标准的书面总结格式  \n无可用内容  \n\n5. 只需要总结内容不要扩展即使内容很少的情况也不要扩展  \n（无有效信息可总结）', '2025-07-05 18:15:21', '2025-07-05 18:15:21', NULL);
INSERT INTO `course` VALUES (20, 39, '数学课', 'static\\media\\2025\\07\\05\\cef33cc655e989740341c7993753d847.mp4', 'static\\media\\2025\\07\\05\\cef33cc655e989740341c7993753d847.mp3', '1. 课程主题  \n电路相关数学问题分析  \n\n2. 主要知识点  \n- 电流、电压、电阻的概念及单位  \n- 欧姆定律公式及变形应用  \n- 串联电路电流、电压规律  \n- 并联电路电流、电压规律  \n- 简单电路计算（单一电源情境）  \n\n3. 难点提醒  \n- 串联与并联电路的电流、电压分配混淆  \n- 欧姆定律公式中物理量对应关系错误  \n- 多步骤计算时单位未统一导致结果偏差  \n- 复杂电路化简为单一回路的方法选择', '2025-07-05 19:10:53', '2025-07-05 19:10:53', NULL);
INSERT INTO `course` VALUES (21, 41, '测试', 'static\\media\\2025\\07\\05\\3fa3694973a1ececa8e6a00d0a56100f.mp4', 'static\\media\\2025\\07\\05\\3fa3694973a1ececa8e6a00d0a56100f.mp3', '1. 课程主题  \n字母认读与发音基础  \n\n2. 主要知识点  \n- 字母A-G的顺序认读  \n- 元音字母A的发音规则（如/æ/）  \n- 辅音字母B、C、D、F、G的标准发音  \n\n3. 难点提醒  \n- 字母B与D的发音混淆  \n- 字母G的发音/g/在词尾的弱化现象  \n- 字母顺序中E-F-G的连贯性记忆  \n\n4. 书面总结格式  \n【课程主题】字母认读与发音基础  \n【核心内容】  \n1. 字母A-G的顺序识记  \n2. 元音A及辅音B/C/D/F/G的标准发音  \n3. 易错点：B/D发音区分、G的发音稳定性  \n\n5. 内容限制说明  \n仅提炼课程显性信息，未补充隐性教学目标或拓展内容。', '2025-07-05 19:26:06', '2025-07-05 19:26:06', NULL);
INSERT INTO `course` VALUES (22, 40, '数学课', 'static\\media\\2025\\07\\05\\a8c9c58993ab1e93d6a6c74ed8ae1781.mp4', 'static\\media\\2025\\07\\05\\a8c9c58993ab1e93d6a6c74ed8ae1781.mp3', '1. **课程主题**  \n   数学课（与电相关的内容）  \n\n2. **主要知识点**  \n   - 电流、电压、电阻的基本概念  \n   - 欧姆定律的公式及应用  \n   - 电路图中的符号识别与分析  \n\n3. **难点提醒**  \n   - 串并联电路中电流与电压的计算区别  \n   - 单位换算（如毫安与安培、千伏与伏特）  \n   - 实际问题中公式的正确选择与变形', '2025-07-05 19:28:35', '2025-07-05 19:28:35', NULL);
INSERT INTO `course` VALUES (23, 21, '音乐课', 'static\\media\\2025\\07\\06\\9caca822917a1001d4250713071f924e.mp4', 'static\\media\\2025\\07\\06\\9caca822917a1001d4250713071f924e.mp3', '1. 课程主题  \n音乐课：基础问候与人工智能认知  \n\n2. 主要知识点  \n（1）音乐课堂问候语“你好”的规范使用  \n（2）人工智能身份识别与简单互动  \n\n3. 难点提醒  \n（1）区分日常口语化问候与课堂仪式化问候的差异  \n（2）理解人工智能回复逻辑与人类教师反馈的区别  \n\n4. 书面总结  \n本节课围绕音乐课堂问候语展开，重点学习“你好”的规范化表达，并通过问答形式建立对人工智能基础认知。需注意课堂问候的正式性要求，以及人工智能与真人教师在互动模式上的本质差异。', '2025-07-06 09:37:45', '2025-07-06 09:37:45', NULL);
INSERT INTO `course` VALUES (24, 21, '音乐课', 'static\\media\\2025\\07\\06\\35b568ec281bb571392dbc4901088742.mp4', 'static\\media\\2025\\07\\06\\35b568ec281bb571392dbc4901088742.mp3', '1. 课程主题  \n音乐课中的人声与乐器认知  \n\n2. 主要知识点  \n- 人声与乐器的音色差异  \n- 简单声音识别（如区分说唱与人声）  \n- 基础音乐互动形式  \n\n3. 难点提醒  \n- 相似音色（如电子合成音与人声）的混淆  \n- 复杂背景音乐中的目标声音提取  \n\n4. 书面总结  \n本节课通过声音示例对比，引导学生初步感知人声与乐器的音色特征，尝试在简单情境中辨别不同声源类型。重点在于建立对声音属性的基础认知，需注意实际听觉环境中的干扰因素对判断的影响。', '2025-07-06 09:40:33', '2025-07-06 09:40:33', NULL);
INSERT INTO `course` VALUES (25, 39, '数学课', 'static\\media\\2025\\07\\06\\5a74dd6977a03b82b2fd314f618930d4.mp4', 'static\\media\\2025\\07\\06\\5a74dd6977a03b82b2fd314f618930d4.mp3', '1. 课程主题  \n简单加法应用（三人总路程计算）  \n\n2. 主要知识点  \n- 加法运算的实际应用  \n- 理解题意中参与人数与路程的关系  \n\n3. 难点提醒  \n- 混淆“共同行走”与“各自行走”的路程计算  \n- 忽略人数导致漏算或重复计算  \n\n4. 总结  \n题目中三人共同行走500米，总路程为500米×3=1500米。需明确“他们”包含三人，避免人数误判。', '2025-07-06 09:46:39', '2025-07-06 09:46:39', NULL);
INSERT INTO `course` VALUES (26, 21, '音乐课', 'static\\media\\2025\\07\\06\\07a92377273f3749a7e73980ecd3fc79.mp4', 'static\\media\\2025\\07\\06\\07a92377273f3749a7e73980ecd3fc79.mp3', '短内容或无法听到你的声音', '2025-07-06 09:49:33', '2025-07-06 09:49:33', NULL);
INSERT INTO `course` VALUES (27, 21, '音乐课', 'static\\media\\2025\\07\\06\\04b604d884993a91318e7143bef05e1f.mp4', 'static\\media\\2025\\07\\06\\04b604d884993a91318e7143bef05e1f.mp3', 'jadslkj', '2025-07-06 09:54:51', '2025-07-06 09:54:51', NULL);
INSERT INTO `course` VALUES (28, 21, '音乐课', 'static\\media\\2025\\07\\06\\4549093f2d0f174097ac3a1103e8d0bd.mp4', 'static\\media\\2025\\07\\06\\4549093f2d0f174097ac3a1103e8d0bd.mp3', 'jadslkj', '2025-07-06 09:57:50', '2025-07-06 09:57:50', NULL);
INSERT INTO `course` VALUES (29, 22, '数学课', 'static\\media\\2025\\07\\06\\553ed55ba2db280d1b370036b3c5840a.mp4', 'static\\media\\2025\\07\\06\\553ed55ba2db280d1b370036b3c5840a.mp3', 'jadslkj', '2025-07-06 09:59:40', '2025-07-06 09:59:40', NULL);
INSERT INTO `course` VALUES (31, 21, '音乐课', 'static\\media\\2025\\07\\06\\25062f967ce5360cf97fd23dc2212fb6.mp4', 'static\\media\\2025\\07\\06\\25062f967ce5360cf97fd23dc2212fb6.mp3', '你好你好，你是人工智能吗？', '2025-07-06 10:07:55', '2025-07-06 10:07:55', NULL);
INSERT INTO `course` VALUES (32, 21, '音乐课', 'static\\media\\2025\\07\\06\\50573fceccf711d39c14f028c8599a7e.mp4', 'static\\media\\2025\\07\\06\\50573fceccf711d39c14f028c8599a7e.mp3', '无', '2025-07-06 10:10:01', '2025-07-06 10:10:01', NULL);
INSERT INTO `course` VALUES (33, 21, '音乐课', 'static\\media\\2025\\07\\06\\33fb25b5c03ded534abd1f549812879d.mp4', 'static\\media\\2025\\07\\06\\33fb25b5c03ded534abd1f549812879d.mp3', '总结：无发听到你的声音', '2025-07-06 10:14:45', '2025-07-06 10:14:45', NULL);
INSERT INTO `course` VALUES (34, 21, '音乐课', 'static\\media\\2025\\07\\06\\c336f92d40d0b3ba5b5bc6083bb29cad.mp4', 'static\\media\\2025\\07\\06\\c336f92d40d0b3ba5b5bc6083bb29cad.mp3', '总结：音乐课 你好你好，你是人工智能吗？', '2025-07-06 10:17:53', '2025-07-06 10:17:53', NULL);
INSERT INTO `course` VALUES (36, 39, 'nnnnnnnnn', 'static\\media\\2025\\07\\06\\8347cacc22662e013886b5e6770385d3.mp4', 'static\\media\\2025\\07\\06\\8347cacc22662e013886b5e6770385d3.mp3', '总结：无发听到你的声音', '2025-07-06 10:28:33', '2025-07-06 10:28:33', NULL);
INSERT INTO `course` VALUES (37, 41, '测试', 'static\\media\\2025\\07\\06\\2352135f1e5bb73cbbcd554db94f8021.mp4', 'static\\media\\2025\\07\\06\\2352135f1e5bb73cbbcd554db94f8021.mp3', '总结：课程内容', '2025-07-06 10:38:21', '2025-07-06 10:38:21', NULL);
INSERT INTO `course` VALUES (41, 39, '体育课', 'static\\media\\2025\\07\\06\\8579c9bae027f6ff45a9754144fb5088.mp4', 'static\\media\\2025\\07\\06\\8579c9bae027f6ff45a9754144fb5088.mp3', '1. 课程主题：体育课  \n2. 主要知识点：无  \n3. 难点提醒：无  \n4. 总结：课程内容', '2025-07-06 11:40:05', '2025-07-06 11:40:05', NULL);
INSERT INTO `course` VALUES (54, 40, '数学课', 'static\\media\\2025\\07\\06\\12640e1d36e9b1d86de38f9c1540a521.mp4', 'static\\media\\2025\\07\\06\\12640e1d36e9b1d86de38f9c1540a521.mp3', '1. **课程主题**  \n1至6的连续自然数加法运算。\n\n2. **主要知识点**  \n- 数字1-6的依次累加。  \n- 加法运算的基本顺序。  \n\n3. **难点提醒**  \n- 连续加法过程中遗漏某个数字。  \n- 计算顺序错误导致结果偏差。  \n\n4. **书面总结**  \n本节课学习了1至6的连续自然数加法运算，重点掌握逐项相加的方法及运算顺序，需注意避免遗漏数字或计算顺序错误。', '2025-07-06 15:04:11', '2025-07-06 15:04:11', NULL);
INSERT INTO `course` VALUES (55, 40, '测试', 'static\\media\\2025\\07\\06\\94b31ca06e2d4011c4df2b8e112d8cbf.mp4', 'static\\media\\2025\\07\\06\\94b31ca06e2d4011c4df2b8e112d8cbf.mp3', '1. **课程主题**  \n连续自然数求和的计算方法。\n\n2. **主要知识点**  \n1. 连续自然数求和公式的应用（如 \\(1+2+\\dots+n = \\frac{n(n+1)}{2}\\)）。  \n2. 等差数列求和公式的理解与简化计算。  \n3. 数学运算中的规律性总结（如首尾配对法）。  \n\n3. **难点提醒**  \n1. 首尾配对时遗漏中间项（针对奇数个项的情况）。  \n2. 公式应用中“项数”与“末项”的混淆（如误将 \\(n\\) 代入错误数值）。  \n3. 计算过程中的乘法或除法步骤出错（如 \\(10 \\times 11 \\div 2\\) 的运算准确性）。', '2025-07-06 15:05:26', '2025-07-06 15:05:26', NULL);
INSERT INTO `course` VALUES (56, 41, '44564', 'static\\media\\2025\\07\\06\\792c0d2c63bccb3ecb7c173e2af8fb59.mp4', 'static\\media\\2025\\07\\06\\792c0d2c63bccb3ecb7c173e2af8fb59.mp3', '1. 课程主题  \n100以内加减法运算规则  \n\n2. 主要知识点  \n- 两位数加法的进位规则  \n- 两位数减法的退位规则  \n- 竖式计算的数位对齐方法  \n- 连加、连减与加减混合运算顺序  \n\n3. 难点提醒  \n- 进位加法与退位减法的步骤混淆  \n- 竖式计算中“0”占位的遗漏  \n- 连加连减时运算顺序的颠倒  \n\n4. 书面总结  \n本节课重点讲解100以内加减法的运算规则，包含进位加法、退位减法、竖式计算及连加连减混合运算。需掌握数位对齐、进位标记、退位借位等核心步骤，注意区分加减法不同场景下的计算逻辑，避免因步骤混淆导致结果错误。', '2025-07-06 15:21:23', '2025-07-06 15:21:23', NULL);
INSERT INTO `course` VALUES (74, 40, '测试', 'static\\media\\2025\\07\\06\\c1d36f265f1a3e8df61f242854fe83f3.mp4', 'static\\media\\2025\\07\\06\\c1d36f265f1a3e8df61f242854fe83f3.mp3', '1. 课程主题  \n页面跳转功能与自动结束机制的基础认知  \n\n2. 主要知识点  \n- 页面跳转的可能性及操作路径  \n- 不同页面功能的基本区分  \n- 自动跳转触发条件的存在性  \n- 课后环节的自动结束设置  \n\n3. 难点提醒  \n- 页面跳转条件与操作结果的对应关系  \n- 自动结束机制触发场景的判断', '2025-07-06 17:10:43', '2025-07-06 17:10:43', NULL);
INSERT INTO `course` VALUES (75, 39, '测试', 'static\\media\\2025\\07\\06\\d021e6e37ff741bc4964c508a5e17466.mp4', 'static\\media\\2025\\07\\06\\d021e6e37ff741bc4964c508a5e17466.mp3', '1. 课程主题  \n工具栏录音功能的基本操作  \n\n2. 主要知识点  \n- 工具栏的位置与组成  \n- 录音功能的启动步骤  \n- 录音文件的保存路径  \n- 操作过程中的常见错误类型  \n\n3. 难点提醒  \n- 工具栏按钮功能的区分（如录音与暂停的图标差异）  \n- 保存路径选择与文件检索的关联性  \n- 操作顺序对录音效果的影响（如未先选择保存位置直接录制）  \n\n书面总结：  \n本次课程围绕工具栏录音功能展开，重点讲解工具栏的位置、录音操作流程及文件管理方法。需注意按钮功能辨识、保存路径设置与操作逻辑顺序，避免因步骤错误导致录音失败或文件丢失。', '2025-07-06 17:14:27', '2025-07-06 17:14:27', NULL);
INSERT INTO `course` VALUES (77, 39, '音乐课', 'static\\media\\2025\\07\\06\\bcfc1e72b548522f6240324314f3fae0.mp4', 'static\\media\\2025\\07\\06\\bcfc1e72b548522f6240324314f3fae0.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 17:33:07', '2025-07-06 17:33:07', NULL);
INSERT INTO `course` VALUES (78, 40, '数学课', 'static\\media\\2025\\07\\06\\30a6708b29ffbe2a83c7ce2d41d96dfe.mp4', 'static\\media\\2025\\07\\06\\30a6708b29ffbe2a83c7ce2d41d96dfe.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 17:36:05', '2025-07-06 17:36:05', NULL);
INSERT INTO `course` VALUES (79, 40, '数据结构', 'static\\media\\2025\\07\\06\\5afcfe7ccaa945245d9e54555de869b7.mp4', 'static\\media\\2025\\07\\06\\5afcfe7ccaa945245d9e54555de869b7.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 17:38:26', '2025-07-06 17:38:26', NULL);
INSERT INTO `course` VALUES (80, 40, '数据结构', 'static\\media\\2025\\07\\06\\1a2d124c10a482386579aaa7125ca4b9.mp4', 'static\\media\\2025\\07\\06\\1a2d124c10a482386579aaa7125ca4b9.mp3', '课程要点：  \n- 课堂涉及数据结构相关概念，提及数论图论、线性表链表及广义表。  \n\n主要知识点：  \n- 数论图论  \n- 线性表链表  \n- 广义表  \n\n（注：文本中夹杂大量非课程内容，如歌词和无关描述，已按规则过滤，仅保留与数据结构直接相关的要点。）', '2025-07-06 17:39:37', '2025-07-06 17:39:37', NULL);
INSERT INTO `course` VALUES (81, 39, 'gfdgfdg', 'static\\media\\2025\\07\\06\\de9e14cff881669f7242bf7efd9451a0.mp4', 'static\\media\\2025\\07\\06\\de9e14cff881669f7242bf7efd9451a0.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 18:33:49', '2025-07-06 18:33:49', NULL);
INSERT INTO `course` VALUES (82, 39, '数据结构', 'static\\media\\2025\\07\\06\\a85c2f9df431cedfd9de8646b9488dd1.mp4', 'static\\media\\2025\\07\\06\\a85c2f9df431cedfd9de8646b9488dd1.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 21:17:37', '2025-07-06 21:17:37', NULL);
INSERT INTO `course` VALUES (83, 21, '116516', 'static\\media\\2025\\07\\06\\da33a8b98f78c9d3641458277a076c8f.mp4', 'static\\media\\2025\\07\\06\\da33a8b98f78c9d3641458277a076c8f.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 21:57:08', '2025-07-06 21:57:08', NULL);
INSERT INTO `course` VALUES (85, 21, '结合客户尽快', 'static\\media\\2025\\07\\06\\31b384dab9f535ab4b7daa767203645c.mp4', 'static\\media\\2025\\07\\06\\31b384dab9f535ab4b7daa767203645c.mp3', '暂无可总结的内容，请延长课程时长再来试试吧', '2025-07-06 23:07:17', '2025-07-06 23:07:17', NULL);

SET FOREIGN_KEY_CHECKS = 1;
