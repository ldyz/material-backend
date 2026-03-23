-- PostgreSQL Database Backup
-- Generated: 2026-02-04 14:21:33
-- Database: materials

-- Disabling triggers and foreign keys for faster import
SET session_replication_role = 'replica';


-- Table: ai_analysis_logs

-- Schema and data for table: ai_analysis_logs

-- Table: construction_log

-- Schema and data for table: construction_log
INSERT INTO "construction_log" ("id", "title", "content", "images", "weather", "project_id", "creator_id", "created_at", "updated_at", "temperature", "progress", "issues", "log_date", "remark") VALUES
(1, '2025-07-02 施工日志 大风 35.72℃', '<p><img src="/static/uploads/construction_log/54abf02213390a1eac484e7f2a129765.jpg"><img src="/static/uploads/construction_log/bfa8b7376e9991763285ff3e63767db6.jpg">1.电缆到货 </p><p>2.电缆保护管和流体用不锈钢管到货</p>', NULL, '晴（白天） 29.51℃', 1, 1, '2025-07-02 00:00:00', '2025-07-30 00:00:00', NULL, NULL, NULL, NULL, NULL),
(2, '2025-07-11 施工日志 多云（白天） 33.31℃', '<p>昨日到货一批材料，清单如下</p><p><img src="/static/uploads/construction_log/19b7509cf6b9d58f848c19d7f96b0df9.png"><img src="/static/uploads/construction_log/a0ee7dfb0a49e1fce2bdd707f17f8a95.png"><img src="/static/uploads/construction_log/5618b07447e3f346bb8b80fdd3c30028.png"><img src="/static/uploads/construction_log/2ee040d6fdc451cb25e837c0321f8296.png"><img src="/static/uploads/construction_log/0167d14b25381d69c6e4a3abfa176928.png"><img src="/static/uploads/construction_log/8ce68fb67022393b11f2de9523231f09.png"><img src="/static/uploads/construction_log/bfae56c6e0f7a5292e1223164ee09588.png"><img src="/static/uploads/construction_log/5f32aa42c1ecec00b4d65a788a2a1a5d.png"><img src="/static/uploads/construction_log/66a125637943c41b52e877f9bb535a44.png"></p>', NULL, '多云（白天） 33.31℃', 1, 1, '2025-07-11 00:00:00', '2025-07-11 00:00:00', NULL, NULL, NULL, NULL, NULL),
(3, '2025-07-16 施工日志 阴 24.81℃', '<p>今日到货控制系统机柜1批，由葛继伟直接保管</p>', NULL, '阴 24.81℃', 1, 1, '2025-07-17 00:00:00', '2025-07-17 00:00:00', NULL, NULL, NULL, NULL, NULL),
(4, '2025-07-21 施工日志 晴（白天） 30.01℃', '<p>1.到货转子流量计1台。</p><p>2.施工人员入场手续已经办理完成，工机具报验已经完成。</p><p>3.今日仪表槽盒拉入现场进行材料入场报验。<img src="/static/uploads/construction_log/7a054b8e89508e260e4acdb73866a465.jpg"></p>', NULL, '晴（白天） 30.01℃', 1, 1, '2025-07-22 00:00:00', '2025-07-22 00:00:00', NULL, NULL, NULL, NULL, NULL),
(5, '2025-07-22 施工日志 晴（白天） 31.31℃', '<p>今日收到热电阻1箱 20支 ；</p><p>明细如下图：</p><p><img src="/static/uploads/construction_log/9e94a6abd66e7c4f979c6fcffcbbb9ac.jpg"><img src="/static/uploads/construction_log/b292cef6ac9cebf926971a4141378350.jpg"></p>', NULL, '晴（白天） 31.31℃', 1, 1, '2025-07-22 00:00:00', '2025-07-22 00:00:00', NULL, NULL, NULL, NULL, NULL),
(6, '2025-07-23 施工日志 晴（白天） 32.71℃', '<p>1、电缆槽盒入场报验</p>', NULL, '晴（白天） 32.71℃', 1, 1, '2025-07-23 00:00:00', '2025-07-23 00:00:00', NULL, NULL, NULL, NULL, NULL),
(7, '2025-07-25 施工日志 多云（白天） 30.01℃', '<p>1.今日到货39台压力、差压变松器<img src="/static/uploads/construction_log/7cde8edd90037027350863d8f3d19386.jpg"><img src="/static/uploads/construction_log/b9c40c9ae59b7032800d933d8a749aab.jpg"></p><p><br></p>', NULL, '多云（白天） 30.01℃', 1, 1, '2025-07-25 00:00:00', '2025-07-25 00:00:00', NULL, NULL, NULL, NULL, NULL),
(8, '2025-07-29 施工日志 多云（白天） 27.51℃', '<p>今日到货扩音对讲系统设备、电视监控系统设备、仪表调节阀13台，双金属温度计19台，压力表26台；</p><p><img src="/static/uploads/construction_log/8432dac77a578d39df8728fb74cdf460.png"><img src="/static/uploads/construction_log/c7b9009735656aebf112914ba98deaf8.png"><img src="/static/uploads/construction_log/259d458cac1cdae6f94d59030d5400bf.png"></p><p><img src="/static/uploads/construction_log/8bb6557677231a08820e4356d29535fe.png"><img src="/static/uploads/construction_log/bd1b84a5447574cf8d51ba5b1db71057.png"><img src="/static/uploads/construction_log/fd97ff3af113f12d84240597ab11d5b9.png"><img src="/static/uploads/construction_log/bf55caa9fb4e34cc10555aa9a5199e59.png"></p><p><br></p>', NULL, '多云（白天） 27.21℃', 1, 1, '2025-07-30 00:00:00', '2025-07-30 00:00:00', NULL, NULL, NULL, NULL, NULL),
(9, '2025-07-30 施工日志 多云（白天） 27.21℃', '<p>今日到货消防系统设备，门禁系统设备，中控控制柜子4套，中控柜子有箱体破损<img src="/static/uploads/construction_log/74a598088c68a1aa59d86b78a2fb8b8a.png"><img src="/static/uploads/construction_log/2dc51744190a92cc29c62c6f88442cfb.png"><img src="/static/uploads/construction_log/9153bca638cb2a15fee46e97a9c67601.png"></p><p><br></p>', NULL, '多云（白天） 27.21℃', 1, 1, '2025-07-30 00:00:00', '2025-07-30 00:00:00', NULL, NULL, NULL, NULL, NULL),
(10, '2025-07-31 施工日志 晴（白天） 30.11℃', '<p>1.今日上午设备联合验收，双金属温度计、压力表、压力变松器、热电阻、火灾报警系统、视频监控系统、门禁系统、广播系统、金属转子流量计等全部验收完毕</p><p>2.今日下午热电阻、双金属温度计、压力表 已经送到检测公司</p><p>3.今日下午13台阀门送到阀门分厂打压</p><p>分包从库房取走如下：</p><p>6分的u型卡子领了300</p><p>一寸的u型卡子领了200</p><p>40的u型卡子领了100</p><p>6分的活接领了15</p><p>一寸的活接领了20</p><p>40的穿板接头领了10</p><p>一寸的穿板接头领了30</p><p>100×100的膨胀螺栓200</p><p>6分的穿板夹片30套</p><p>6分穿板接头没有用夹片代替</p><p>蛇皮软管3卷 14+12+34=60米</p><p>摄像系统：摄像头1个，支架1个，非防爆控制箱2套</p>', NULL, '晴（白天） 20℃~29.11℃', 1, 1, '2025-07-31 00:00:00', '2025-08-04 00:00:00', NULL, NULL, NULL, NULL, NULL),
(11, '2025-08-04 施工日志 晴（白天） 20℃~29.24℃', '<p>1、今日到货紧固件1箱</p><p>2、到货e+h电导率测试仪1箱</p><p><br></p>', NULL, '晴（白天） 20℃~29.24℃', 1, 1, '2025-08-04 00:00:00', '2025-08-04 00:00:00', NULL, NULL, NULL, NULL, NULL),
(12, '2025-08-05 施工日志 多云（白天） 19.9℃~28.51℃', '<p>1.到货涡街流量计4台</p><p><img src="/static/uploads/construction_log/08a27ea75b222d2e5f9bab15dcc322e7.jpg"></p><p>2、今日将所有未到材料标注后返给总包了</p>', NULL, '多云（白天） 19.9℃~30.41℃', 1, 1, '2025-08-05 00:00:00', '2025-08-05 00:00:00', NULL, NULL, NULL, NULL, NULL),
(13, '2025-08-07 施工日志 多云（白天） 18.51℃~27.31℃', '<p>1.今日到货压力表2块 位号8010 8011,双金属温度计2块 位号8010 8011</p><p>2.缠绕垫片等到货，在库房里侧窗台上</p><p><br></p><p><br></p>', NULL, '多云（白天） 18.51℃~27.31℃', 1, 1, '2025-08-07 00:00:00', '2025-08-07 00:00:00', NULL, NULL, NULL, NULL, NULL),
(14, '2025-08-20 施工日志 晴 17℃~29℃', '<p>1.电磁流量计7台 送检</p>', NULL, '定位失败', 1, 1, '2025-08-22 00:00:00', '2025-08-22 00:00:00', NULL, NULL, NULL, NULL, NULL),
(15, '2025-08-21施工日志 多云 17℃~27.1℃', '<p>1.电信设备全部取走</p><p>2.3台火灾报警模块箱取走</p>', NULL, '定位失败', 1, 1, '2025-08-22 00:00:00', '2025-08-22 00:00:00', NULL, NULL, NULL, NULL, NULL),
(16, '2025-09-08 施工日志 晴（白天） 14℃~20.18℃', '<p>迪雷来取资料：</p><p>1.压力表校验记录 </p><p>2.广播系统所有合格证和质量证明文件</p><p>3.视频监控系统所有合格证和质量证明文件</p><p>4.火灾系统所有合格证和质量证明文件</p><p>5.调节阀合格证和说明书</p><p>6.双金属、热电阻、压力表合格证 所有</p><p>7.金属缠绕垫质量证明文件</p><p>8.分包小宋借走电缆槽盒100*50 12节 24米，包含槽盒螺丝和连接板<img src="/static/uploads/construction_log/8e2dcd23049b1c766d8dc8ecca6e87ba.jpg"><img src="/static/uploads/construction_log/3beba799a4cbc8ef6337c5d90d568752.jpg"><img src="/static/uploads/construction_log/bebed7df539ee5fe35c17a072caef73d.jpg"><img src="/static/uploads/construction_log/20dfa75b0c6aac92bbb335ea3b4eb66c.jpg"><img src="/static/uploads/construction_log/ea35079c83004f2d471352f63a203395.jpg"><img src="/static/uploads/construction_log/c332593c7ce72aa064cfd9c8d2b7c4a7.jpg"><img src="/static/uploads/construction_log/d1c08968dd970187ba827ea2ef5928e2.jpg"><img src="/static/uploads/construction_log/00292c274517759ed05e80e19633cd97.jpg"><img src="/static/uploads/construction_log/e96bb29fbbb13509bf98cf44ff6b1e96.jpg"><img src="/static/uploads/construction_log/4672396a3d6f406f950c0453855fb0c4.jpg"><img src="/static/uploads/construction_log/a230fca0a11b2274e1491282f67c5779.jpg"><img src="/static/uploads/construction_log/e4ad5f004dc18f7b0de47f79e4c6838d.jpg"></p>', NULL, '晴（白天） 14℃~20.18℃', 1, 1, '2025-09-08 00:00:00', '2025-09-08 00:00:00', NULL, NULL, NULL, NULL, NULL),
(19, '2026年1月30日 施工日志', '<ol><li>施工材料计划提报完成，电缆已经确定排产；年后发货；</li><li>山东科技斐林试剂l</li><li>塑料袋放假快乐； j；</li><li>乐山大佛进口量；是</li><li>是；代理费会计师； l</li><li>塑料袋放假啊； <img src="/static/uploads/1770179465592_5592.jpg"></li><li>速度；立方空间； l</li></ol><p><br></p>', '', 'sunny', 10, 1, '2026-01-30 00:00:00', '2026-02-04 00:00:00', -12, '', '', '2026-02-04', '');

-- Table: inbound_items

-- Schema and data for table: inbound_items
INSERT INTO "inbound_items" ("id", "inbound_order_id", "stock_id", "material_id", "quantity", "unit_price", "status", "remark", "created_at") VALUES
(1, 1, NULL, 2, '40.000', '0.00', 'pending', '', '2026-02-03 08:52:09'),
(2, 1, NULL, 4, '16.000', '0.00', 'pending', '', '2026-02-03 08:52:09'),
(5, 3, NULL, 6, '260.000', '0.00', 'pending', '', '2026-02-04 05:12:29'),
(6, 3, NULL, 8, '2.000', '0.00', 'pending', '', '2026-02-04 05:12:29'),
(3, 2, NULL, 2, '40.000', '0.00', 'pending', '', '2026-02-04 04:38:01'),
(4, 2, NULL, 4, '16.000', '0.00', 'pending', '', '2026-02-04 04:38:01'),
(7, 4, NULL, 10, '6.000', '0.00', 'pending', '', '2026-02-04 06:36:28'),
(8, 4, NULL, 12, '12.000', '0.00', 'pending', '', '2026-02-04 06:36:28'),
(9, 5, NULL, 14, '198.000', '0.00', 'pending', '', '2026-02-04 06:45:23'),
(10, 5, NULL, 16, '4.000', '0.00', 'pending', '', '2026-02-04 06:45:23'),
(11, 6, NULL, 14, '2.000', '0.00', 'pending', '', '2026-02-04 06:45:52'),
(12, 7, NULL, 18, '2.000', '0.00', 'pending', '', '2026-02-04 06:49:14'),
(13, 8, NULL, 18, '2.000', '0.00', 'pending', '', '2026-02-04 07:54:56'),
(14, 8, NULL, 20, '78.000', '0.00', 'pending', '', '2026-02-04 07:54:56'),
(15, 9, NULL, 22, '13.000', '0.00', 'pending', '', '2026-02-04 08:53:25'),
(16, 10, NULL, 24, '6.000', '0.00', 'pending', '', '2026-02-04 08:58:33'),
(17, 11, NULL, 26, '560.000', '0.00', 'pending', '', '2026-02-04 13:37:23'),
(18, 12, NULL, 28, '100.000', '0.00', 'pending', '', '2026-02-04 13:38:07'),
(19, 13, NULL, 30, '2.000', '0.00', 'pending', '', '2026-02-04 13:39:54'),
(20, 14, NULL, 32, '8.000', '0.00', 'pending', '', '2026-02-04 13:45:03'),
(21, 14, NULL, 34, '8.000', '0.00', 'pending', '', '2026-02-04 13:45:03'),
(22, 15, NULL, 36, '4.000', '0.00', 'pending', '', '2026-02-04 13:50:13'),
(23, 15, NULL, 38, '6.000', '0.00', 'pending', '', '2026-02-04 13:50:13'),
(24, 16, NULL, 42, '12.000', '0.00', 'pending', '', '2026-02-04 13:54:03'),
(25, 17, NULL, 40, '24.000', '0.00', 'pending', '', '2026-02-04 14:04:57');

-- Table: inbound_orders

-- Schema and data for table: inbound_orders
INSERT INTO "inbound_orders" ("id", "order_no", "supplier", "contact", "project_id", "creator_id", "creator_name", "status", "notes", "remark", "total_amount", "created_at", "updated_at", "plan_id") VALUES
(12, 'RK20260204133807', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 13:38:07', '2026-02-04 13:39:09', 1),
(13, 'RK20260204133954', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 13:39:54', '2026-02-04 13:44:35', 1),
(14, 'RK20260204134503', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 13:45:03', '2026-02-04 13:49:54', 1),
(15, 'RK20260204135013', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 13:50:13', '2026-02-04 13:52:54', 1),
(16, 'RK20260204135403', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 13:54:03', '2026-02-04 13:54:07', 1),
(17, 'RK20260204140457', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 14:04:57', '2026-02-04 14:05:01', 1),
(1, 'RK20260203085209', '对对对', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-03 08:52:09', '2026-02-03 08:52:13', 1),
(3, 'RK20260204051229', '我摸', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 05:12:29', '2026-02-04 06:30:59', 1),
(2, 'RK20260204043801', 'fff', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 04:38:01', '2026-02-04 06:31:18', 1),
(4, 'RK20260204063628', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 06:36:28', '2026-02-04 06:36:32', 1),
(5, 'RK20260204064523', '供应商', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 06:45:23', '2026-02-04 06:45:26', 1),
(6, 'RK20260204064552', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 06:45:52', '2026-02-04 06:45:56', 1),
(7, 'RK20260204064914', '测试', '', 10, 1, 'admin', 'rejected', '', '', 0, '2026-02-04 06:49:14', '2026-02-04 06:49:19', 1),
(8, 'RK20260204075456', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 07:54:56', '2026-02-04 07:55:00', 1),
(9, 'RK20260204085325', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 08:53:25', '2026-02-04 08:53:32', 1),
(10, 'RK20260204085833', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 08:58:33', '2026-02-04 08:58:41', 1),
(11, 'RK20260204133723', '测试', '', 10, 1, 'admin', 'completed', '', '', 0, '2026-02-04 13:37:23', '2026-02-04 13:37:32', 1);

-- Table: material_categories

-- Schema and data for table: material_categories
INSERT INTO "material_categories" ("id", "name", "code", "sort", "remark", "created_at", "updated_at", "parent_id", "level", "path") VALUES
(27, '钢材', 'STEEL', 1, '各种钢材材料，如钢筋、钢板、钢管等', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(28, '水泥', 'CEMENT', 2, '各类水泥及水泥制品', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(29, '砂石', 'SAND_STONE', 3, '砂、石、砂石等骨料', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(30, '电气材料', 'ELECTRICAL', 4, '电线、电缆、开关、插座等电气材料', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(31, '管道材料', 'PIPE', 5, '各类管道及管件', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(32, '木材', 'WOOD', 6, '原木、板材、木方等木材', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(33, '涂料', 'PAINT', 7, '各类油漆、涂料', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(34, '保温材料', 'INSULATION', 8, '保温棉、泡沫板等保温材料', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(35, '防水材料', 'WATERPROOF', 9, '防水卷材、防水涂料等', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(36, '五金配件', 'HARDWARE', 10, '螺丝、螺母、钉子等五金配件', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(37, '劳保用品', 'SAFETY', 11, '安全帽、手套、工作服等劳保用品', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(38, '工具', 'TOOLS', 12, '电动工具、手动工具等', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL),
(39, '其他', 'OTHER', 99, '其他物资', '2026-02-02 05:33:09', '2026-02-02 05:33:09', 0, 1, NULL);

-- Table: material_master

-- Schema and data for table: material_master
INSERT INTO "material_master" ("id", "code", "name", "specification", "unit", "category", "safety_stock", "description", "created_at", "updated_at") VALUES
(2, 'AUTO1', '螺母', 'M16', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(4, 'AUTO3', '锁紧螺母', 'G3/4"', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(6, 'AUTO5', '阻燃电缆', 'ZA-DJYPVRP 4*2*1.5', '米', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(8, 'AUTO7', '金属缠绕式垫片', 'DN20 PN20 带内环和对中环型', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(10, 'AUTO9', '防爆挠性连接管', '1/2"NPT(外）-G3/4"(内)  Exdb IIC T4    IP65  L=700mm', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(12, 'AUTO11', '等边角钢', '∠50*50*5', '米', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(14, 'AUTO13', '阻燃本安电缆', 'ZA-ia-DJYVRP 1*2*1.5', '米', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(16, 'AUTO15', '金属缠绕式垫片', 'DN40 PN20 带内环和对中环型', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(18, 'AUTO17', '金属缠绕式垫片', 'DN80 PN20 带内环和对中环型', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(20, 'AUTO19', '镀锌焊接钢管', 'DN20  φ26.9*2.8', '米', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(22, 'AUTO21', '带盖弯通', 'G3/4"', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(24, 'AUTO23', '护线帽', 'G3/4"', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(26, 'AUTO25', '阻燃控制电缆', 'ZA-DJYVRP 1*2*1.5', '米', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(28, 'AUTO27', '阻燃控制电缆', 'ZA-KVVRP  2*2.5', '米', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(30, 'AUTO29', '带颈对焊法兰', 'DN20 PN20 WN/RF', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(32, 'AUTO31', '全螺纹螺柱', 'M14*75', '条', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(34, 'AUTO33', '全螺纹螺柱', 'M16*135', '条', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(36, 'AUTO35', '全螺纹螺柱', 'M14*85', '条', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(38, 'AUTO37', '防爆非铠装电缆密封接头', '1/2"NPT(外）-1/2"NPT(内)(电缆外径13.5) Exdb IIC T4 Gb IP65', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(40, 'AUTO39', '螺母', 'M14', '个', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(42, 'AUTO41', '全螺纹螺柱', 'M16*85', '条', '', '0.000', NULL, '2026-02-03 08:51:22', '2026-02-03 08:51:22');

-- Table: material_plan_items

-- Schema and data for table: material_plan_items
INSERT INTO "material_plan_items" ("id", "plan_id", "material_id", "planned_quantity", "unit_price", "required_date", "priority", "status", "remark", "created_at", "updated_at") VALUES
(1, 1, 2, '40.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(2, 1, 4, '16.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(3, 1, 6, '260.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(4, 1, 8, '2.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(5, 1, 10, '6.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(6, 1, 12, '12.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(7, 1, 14, '200.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(8, 1, 16, '4.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(9, 1, 18, '2.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(10, 1, 20, '78.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(11, 1, 22, '13.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(12, 1, 24, '6.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(13, 1, 26, '560.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(14, 1, 28, '100.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(15, 1, 30, '2.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(16, 1, 32, '8.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(17, 1, 34, '8.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(18, 1, 36, '4.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(19, 1, 38, '6.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(20, 1, 40, '24.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22'),
(21, 1, 42, '12.000', '0.00', NULL, 'normal', 'pending', '', '2026-02-03 08:51:22', '2026-02-03 08:51:22');

-- Table: material_plans

-- Schema and data for table: material_plans
INSERT INTO "material_plans" ("id", "plan_no", "plan_name", "project_id", "plan_type", "status", "priority", "planned_start_date", "planned_end_date", "total_budget", "actual_cost", "description", "remark", "workflow_instance_id", "creator_id", "creator_name", "approver_id", "approver_name", "approved_at", "created_at", "updated_at") VALUES
(1, 'MP2602030001', '计划', 10, 'procurement', 'approved', 'normal', NULL, NULL, 0, 0, '', '', 1, 1, 'admin', 1, 'admin', '2026-02-03 08:51:43', '2026-02-03 08:51:22', '2026-02-03 08:51:43');

-- Table: notifications

-- Schema and data for table: notifications
INSERT INTO "notifications" ("id", "user_id", "type", "title", "content", "data", "is_read", "created_at", "read_at") VALUES
(94, 1, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 32, "business_no": "MP2602020001", "instance_id": 34, "business_type": "material_plan"}', false, '2026-02-02 04:10:30', NULL),
(95, 3, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 32, "business_no": "MP2602020001", "instance_id": 34, "business_type": "material_plan"}', false, '2026-02-02 04:10:30', NULL),
(96, 4, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 32, "business_no": "MP2602020001", "instance_id": 34, "business_type": "material_plan"}', false, '2026-02-02 04:10:30', NULL),
(97, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 39, "business_no": "CK20260202001", "instance_id": 35, "business_type": "requisition"}', false, '2026-02-02 04:11:31', NULL),
(98, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 39, "business_no": "CK20260202001", "instance_id": 35, "business_type": "requisition"}', false, '2026-02-02 04:11:31', NULL),
(99, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 39, "business_no": "CK20260202001", "instance_id": 35, "business_type": "requisition"}', false, '2026-02-02 04:11:31', NULL),
(100, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 39, "requisition_no": "CK20260202001"}', false, '2026-02-02 04:11:31', NULL),
(101, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 39, "requisition_no": "CK20260202001"}', false, '2026-02-02 04:11:31', NULL),
(102, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 39, "requisition_no": "CK20260202001"}', false, '2026-02-02 04:11:31', NULL),
(103, 1, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 33, "business_no": "MP2602020002", "instance_id": 36, "business_type": "material_plan"}', false, '2026-02-02 04:47:29', NULL),
(104, 3, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 33, "business_no": "MP2602020002", "instance_id": 36, "business_type": "material_plan"}', false, '2026-02-02 04:47:29', NULL),
(105, 4, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 33, "business_no": "MP2602020002", "instance_id": 36, "business_type": "material_plan"}', false, '2026-02-02 04:47:29', NULL),
(106, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 40, "business_no": "CK20260202002", "instance_id": 37, "business_type": "requisition"}', false, '2026-02-02 05:24:34', NULL),
(107, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 40, "business_no": "CK20260202002", "instance_id": 37, "business_type": "requisition"}', false, '2026-02-02 05:24:34', NULL),
(108, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 40, "business_no": "CK20260202002", "instance_id": 37, "business_type": "requisition"}', false, '2026-02-02 05:24:34', NULL),
(109, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202002', '{"applicant": "admin", "project_id": 10, "items_count": 2, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 40, "requisition_no": "CK20260202002"}', false, '2026-02-02 05:24:34', NULL),
(110, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202002', '{"applicant": "admin", "project_id": 10, "items_count": 2, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 40, "requisition_no": "CK20260202002"}', false, '2026-02-02 05:24:34', NULL),
(111, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202002', '{"applicant": "admin", "project_id": 10, "items_count": 2, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 40, "requisition_no": "CK20260202002"}', false, '2026-02-02 05:24:34', NULL),
(112, 1, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 34, "business_no": "MP2602020001", "instance_id": 38, "business_type": "material_plan"}', false, '2026-02-02 05:33:39', NULL),
(113, 3, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 34, "business_no": "MP2602020001", "instance_id": 38, "business_type": "material_plan"}', false, '2026-02-02 05:33:39', NULL),
(114, 4, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602020001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 34, "business_no": "MP2602020001", "instance_id": 38, "business_type": "material_plan"}', false, '2026-02-02 05:33:39', NULL),
(115, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 41, "business_no": "CK20260202001", "instance_id": 39, "business_type": "requisition"}', false, '2026-02-02 05:34:39', NULL),
(116, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 41, "business_no": "CK20260202001", "instance_id": 39, "business_type": "requisition"}', false, '2026-02-02 05:34:39', NULL),
(117, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 41, "business_no": "CK20260202001", "instance_id": 39, "business_type": "requisition"}', false, '2026-02-02 05:34:39', NULL),
(118, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 41, "requisition_no": "CK20260202001"}', false, '2026-02-02 05:34:39', NULL),
(119, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 41, "requisition_no": "CK20260202001"}', false, '2026-02-02 05:34:39', NULL),
(120, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 41, "requisition_no": "CK20260202001"}', false, '2026-02-02 05:34:39', NULL),
(121, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 42, "business_no": "CK20260202002", "instance_id": 40, "business_type": "requisition"}', false, '2026-02-02 05:40:26', NULL),
(122, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 42, "business_no": "CK20260202002", "instance_id": 40, "business_type": "requisition"}', false, '2026-02-02 05:40:26', NULL),
(123, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260202002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 42, "business_no": "CK20260202002", "instance_id": 40, "business_type": "requisition"}', false, '2026-02-02 05:40:26', NULL),
(124, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202002', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 42, "requisition_no": "CK20260202002"}', false, '2026-02-02 05:40:26', NULL),
(125, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202002', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 42, "requisition_no": "CK20260202002"}', false, '2026-02-02 05:40:26', NULL),
(126, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260202002', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 42, "requisition_no": "CK20260202002"}', false, '2026-02-02 05:40:26', NULL),
(127, 1, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 10, "business_no": "MP2602030001", "instance_id": 41, "business_type": "material_plan"}', false, '2026-02-03 04:20:18', NULL),
(128, 3, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 10, "business_no": "MP2602030001", "instance_id": 41, "business_type": "material_plan"}', false, '2026-02-03 04:20:18', NULL),
(129, 4, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 10, "business_no": "MP2602030001", "instance_id": 41, "business_type": "material_plan"}', false, '2026-02-03 04:20:18', NULL),
(130, 1, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 11, "business_no": "MP2602030002", "instance_id": 42, "business_type": "material_plan"}', false, '2026-02-03 04:22:35', NULL),
(131, 3, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 11, "business_no": "MP2602030002", "instance_id": 42, "business_type": "material_plan"}', false, '2026-02-03 04:22:35', NULL),
(132, 4, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 11, "business_no": "MP2602030002", "instance_id": 42, "business_type": "material_plan"}', false, '2026-02-03 04:22:36', NULL),
(133, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "CK20260203001", "instance_id": 43, "business_type": "requisition"}', false, '2026-02-03 05:40:24', NULL),
(134, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "CK20260203001", "instance_id": 43, "business_type": "requisition"}', false, '2026-02-03 05:40:24', NULL),
(135, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "CK20260203001", "instance_id": 43, "business_type": "requisition"}', false, '2026-02-03 05:40:24', NULL),
(136, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 1, "requisition_no": "CK20260203001"}', false, '2026-02-03 05:40:24', NULL),
(137, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 1, "requisition_no": "CK20260203001"}', false, '2026-02-03 05:40:24', NULL),
(138, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 1, "requisition_no": "CK20260203001"}', false, '2026-02-03 05:40:24', NULL),
(139, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 2, "business_no": "CK20260203002", "instance_id": 44, "business_type": "requisition"}', false, '2026-02-03 06:01:51', NULL),
(140, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 2, "business_no": "CK20260203002", "instance_id": 44, "business_type": "requisition"}', false, '2026-02-03 06:01:51', NULL),
(141, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203002', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 2, "business_no": "CK20260203002", "instance_id": 44, "business_type": "requisition"}', false, '2026-02-03 06:01:51', NULL),
(142, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203002', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 2, "requisition_no": "CK20260203002"}', false, '2026-02-03 06:01:51', NULL),
(143, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203002', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 2, "requisition_no": "CK20260203002"}', false, '2026-02-03 06:01:51', NULL),
(144, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203002', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 2, "requisition_no": "CK20260203002"}', false, '2026-02-03 06:01:51', NULL),
(145, 1, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "MP2602030001", "instance_id": 1, "business_type": "material_plan"}', false, '2026-02-03 08:51:37', NULL),
(146, 3, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "MP2602030001", "instance_id": 1, "business_type": "material_plan"}', false, '2026-02-03 08:51:37', NULL),
(147, 4, 'system', '待审批：业务单据', '业务单据 新节点 需要您审批，单号：MP2602030001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "MP2602030001", "instance_id": 1, "business_type": "material_plan"}', false, '2026-02-03 08:51:37', NULL),
(148, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "CK20260203001", "instance_id": 2, "business_type": "requisition"}', false, '2026-02-03 09:01:27', NULL),
(149, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "CK20260203001", "instance_id": 2, "business_type": "requisition"}', false, '2026-02-03 09:01:27', NULL),
(150, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260203001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 1, "business_no": "CK20260203001", "instance_id": 2, "business_type": "requisition"}', false, '2026-02-03 09:01:27', NULL),
(151, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203001', '{"applicant": "admin", "project_id": 10, "items_count": 2, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 1, "requisition_no": "CK20260203001"}', false, '2026-02-03 09:01:28', NULL),
(152, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203001', '{"applicant": "admin", "project_id": 10, "items_count": 2, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 1, "requisition_no": "CK20260203001"}', false, '2026-02-03 09:01:28', NULL),
(153, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260203001', '{"applicant": "admin", "project_id": 10, "items_count": 2, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 1, "requisition_no": "CK20260203001"}', false, '2026-02-03 09:01:28', NULL),
(154, 1, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260204001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 2, "business_no": "CK20260204001", "instance_id": 3, "business_type": "requisition"}', false, '2026-02-04 05:42:43', NULL),
(155, 3, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260204001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 2, "business_no": "CK20260204001", "instance_id": 3, "business_type": "requisition"}', false, '2026-02-04 05:42:43', NULL),
(156, 4, 'requisition_approve', '待审批：领料单', '领料单 新节点 需要您审批，单号：CK20260204001', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 2, "business_no": "CK20260204001", "instance_id": 3, "business_type": "requisition"}', false, '2026-02-04 05:42:43', NULL),
(157, 1, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260204001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 2, "requisition_no": "CK20260204001"}', false, '2026-02-04 05:42:43', NULL),
(158, 3, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260204001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 2, "requisition_no": "CK20260204001"}', false, '2026-02-04 05:42:43', NULL),
(159, 4, 'requisition_approve', '新的出库单待审批', '申请人：admin，项目：线性装置聚烯烃弹性体生产优化改造，单号：CK20260204001', '{"applicant": "admin", "project_id": 10, "items_count": 1, "project_name": "线性装置聚烯烃弹性体生产优化改造", "requisition_id": 2, "requisition_no": "CK20260204001"}', false, '2026-02-04 05:42:43', NULL),
(160, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133723', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 11, "business_no": "RK20260204133723", "instance_id": 4, "business_type": "inbound_order"}', false, '2026-02-04 13:37:23', NULL),
(161, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133723', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 11, "business_no": "RK20260204133723", "instance_id": 4, "business_type": "inbound_order"}', false, '2026-02-04 13:37:23', NULL),
(162, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133723', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 11, "business_no": "RK20260204133723", "instance_id": 4, "business_type": "inbound_order"}', false, '2026-02-04 13:37:23', NULL),
(163, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133807', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 12, "business_no": "RK20260204133807", "instance_id": 5, "business_type": "inbound_order"}', false, '2026-02-04 13:38:07', NULL),
(164, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133807', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 12, "business_no": "RK20260204133807", "instance_id": 5, "business_type": "inbound_order"}', false, '2026-02-04 13:38:07', NULL),
(165, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133807', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 12, "business_no": "RK20260204133807", "instance_id": 5, "business_type": "inbound_order"}', false, '2026-02-04 13:38:07', NULL),
(166, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133954', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 13, "business_no": "RK20260204133954", "instance_id": 6, "business_type": "inbound_order"}', false, '2026-02-04 13:39:55', NULL),
(167, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133954', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 13, "business_no": "RK20260204133954", "instance_id": 6, "business_type": "inbound_order"}', false, '2026-02-04 13:39:55', NULL),
(168, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204133954', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 13, "business_no": "RK20260204133954", "instance_id": 6, "business_type": "inbound_order"}', false, '2026-02-04 13:39:55', NULL),
(169, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204134503', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 14, "business_no": "RK20260204134503", "instance_id": 7, "business_type": "inbound_order"}', false, '2026-02-04 13:45:03', NULL),
(170, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204134503', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 14, "business_no": "RK20260204134503", "instance_id": 7, "business_type": "inbound_order"}', false, '2026-02-04 13:45:03', NULL),
(171, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204134503', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 14, "business_no": "RK20260204134503", "instance_id": 7, "business_type": "inbound_order"}', false, '2026-02-04 13:45:03', NULL),
(172, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204135013', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 15, "business_no": "RK20260204135013", "instance_id": 8, "business_type": "inbound_order"}', false, '2026-02-04 13:50:13', NULL),
(173, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204135013', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 15, "business_no": "RK20260204135013", "instance_id": 8, "business_type": "inbound_order"}', false, '2026-02-04 13:50:13', NULL),
(174, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204135013', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 15, "business_no": "RK20260204135013", "instance_id": 8, "business_type": "inbound_order"}', false, '2026-02-04 13:50:13', NULL),
(175, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204135403', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 16, "business_no": "RK20260204135403", "instance_id": 9, "business_type": "inbound_order"}', false, '2026-02-04 13:54:03', NULL),
(176, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204135403', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 16, "business_no": "RK20260204135403", "instance_id": 9, "business_type": "inbound_order"}', false, '2026-02-04 13:54:03', NULL),
(177, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204135403', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 16, "business_no": "RK20260204135403", "instance_id": 9, "business_type": "inbound_order"}', false, '2026-02-04 13:54:03', NULL),
(178, 1, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204140457', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 17, "business_no": "RK20260204140457", "instance_id": 10, "business_type": "inbound_order"}', false, '2026-02-04 14:04:58', NULL),
(179, 3, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204140457', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 17, "business_no": "RK20260204140457", "instance_id": 10, "business_type": "inbound_order"}', false, '2026-02-04 14:04:58', NULL),
(180, 4, 'inbound_approve', '待审批：入库单', '入库单 新节点 需要您审批，单号：RK20260204140457', '{"node_key": "node_1", "initiator": "admin", "node_name": "新节点", "business_id": 17, "business_no": "RK20260204140457", "instance_id": 10, "business_type": "inbound_order"}', false, '2026-02-04 14:04:58', NULL);

-- Table: operation_logs

-- Schema and data for table: operation_logs
INSERT INTO "operation_logs" ("id", "user_id", "username", "operation", "module", "resource_type", "resource_id", "resource_no", "changes", "request_method", "request_path", "request_params", "ip_address", "user_agent", "status", "error_message", "created_at", "updated_at") VALUES
(1, 1, 'admin', 'create', 'inbound', 'InboundOrder', 2, 'RK20260204043801', NULL, '', '', '{"items": [{"remark": "", "quantity": 40, "unit_price": 0, "material_id": 2}, {"remark": "", "quantity": 16, "unit_price": 0, "material_id": 4}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "fff", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 04:38:01', '2026-02-04 04:38:01'),
(2, 1, 'admin', 'create', 'inbound', 'InboundOrder', 3, 'RK20260204051229', NULL, '', '', '{"items": [{"remark": "", "quantity": 260, "unit_price": 0, "material_id": 6}, {"remark": "", "quantity": 2, "unit_price": 0, "material_id": 8}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "我摸", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 05:12:29', '2026-02-04 05:12:29'),
(3, 1, 'admin', 'create', 'requisition', 'Requisition', 2, 'CK20260204001', NULL, '', '', '{"items": [{"remark": "", "stock_id": 1, "material_id": 2, "requested_quantity": 5}], "urgent": false, "plan_id": null, "purpose": "测试审核", "applicant": "admin", "department": "", "project_id": 10}', '', '', 'success', '', '2026-02-04 05:42:43', '2026-02-04 05:42:43'),
(4, 1, 'admin', 'approve', 'requisition', 'Requisition', 2, 'CK20260204001', NULL, '', '', '{"comment": "测试审核通过"}', '', '', 'success', '', '2026-02-04 05:42:50', '2026-02-04 05:42:50'),
(5, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 3, 'RK20260204051229', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 06:30:59', '2026-02-04 06:30:59'),
(6, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 2, 'RK20260204043801', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 06:31:18', '2026-02-04 06:31:18'),
(7, 1, 'admin', 'create', 'requisition', 'RequisitionIssue', 2, 'CK20260204001', NULL, '', '', '{"items": null, "notes": "", "remark": ""}', '', '', 'success', '', '2026-02-04 06:36:01', '2026-02-04 06:36:01'),
(8, 1, 'admin', 'create', 'inbound', 'InboundOrder', 4, 'RK20260204063628', NULL, '', '', '{"items": [{"remark": "", "quantity": 6, "unit_price": 0, "material_id": 10}, {"remark": "", "quantity": 12, "unit_price": 0, "material_id": 12}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 06:36:28', '2026-02-04 06:36:28'),
(9, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 4, 'RK20260204063628', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 06:36:32', '2026-02-04 06:36:32'),
(10, 1, 'admin', 'create', 'inbound', 'InboundOrder', 5, 'RK20260204064523', NULL, '', '', '{"items": [{"remark": "", "quantity": 198, "unit_price": 0, "material_id": 14}, {"remark": "", "quantity": 4, "unit_price": 0, "material_id": 16}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "供应商", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 06:45:23', '2026-02-04 06:45:23'),
(11, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 5, 'RK20260204064523', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 06:45:26', '2026-02-04 06:45:26'),
(12, 1, 'admin', 'create', 'inbound', 'InboundOrder', 6, 'RK20260204064552', NULL, '', '', '{"items": [{"remark": "", "quantity": 2, "unit_price": 0, "material_id": 14}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 06:45:52', '2026-02-04 06:45:52'),
(13, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 6, 'RK20260204064552', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 06:45:56', '2026-02-04 06:45:56'),
(14, 1, 'admin', 'create', 'inbound', 'InboundOrder', 7, 'RK20260204064914', NULL, '', '', '{"items": [{"remark": "", "quantity": 2, "unit_price": 0, "material_id": 18}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 06:49:15', '2026-02-04 06:49:15'),
(15, 1, '未知用户', 'reject', 'inbound', 'InboundOrder', 7, 'RK20260204064914', NULL, '', '', '{"reason": ""}', '', '', 'success', '', '2026-02-04 06:49:19', '2026-02-04 06:49:19'),
(16, 1, 'admin', 'create', 'inbound', 'InboundOrder', 8, 'RK20260204075456', NULL, '', '', '{"items": [{"remark": "", "quantity": 2, "unit_price": 0, "material_id": 18}, {"remark": "", "quantity": 78, "unit_price": 0, "material_id": 20}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 07:54:56', '2026-02-04 07:54:56'),
(17, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 8, 'RK20260204075456', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 07:55:00', '2026-02-04 07:55:00'),
(18, 1, 'admin', 'create', 'inbound', 'InboundOrder', 9, 'RK20260204085325', NULL, '', '', '{"items": [{"remark": "", "quantity": 13, "unit_price": 0, "material_id": 22}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 08:53:25', '2026-02-04 08:53:25'),
(19, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 9, 'RK20260204085325', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 08:53:32', '2026-02-04 08:53:32'),
(20, 1, 'admin', 'create', 'inbound', 'InboundOrder', 10, 'RK20260204085833', NULL, '', '', '{"items": [{"remark": "", "quantity": 6, "unit_price": 0, "material_id": 24}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 08:58:33', '2026-02-04 08:58:33'),
(21, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 10, 'RK20260204085833', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 08:58:41', '2026-02-04 08:58:41'),
(22, 1, 'admin', 'create', 'inbound', 'InboundOrder', 11, 'RK20260204133723', NULL, '', '', '{"items": [{"remark": "", "quantity": 560, "unit_price": 0, "material_id": 26}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 13:37:23', '2026-02-04 13:37:23'),
(23, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 11, 'RK20260204133723', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 13:37:32', '2026-02-04 13:37:32'),
(24, 1, 'admin', 'create', 'inbound', 'InboundOrder', 12, 'RK20260204133807', NULL, '', '', '{"items": [{"remark": "", "quantity": 100, "unit_price": 0, "material_id": 28}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 13:38:07', '2026-02-04 13:38:07'),
(25, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 12, 'RK20260204133807', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 13:39:09', '2026-02-04 13:39:09'),
(26, 1, 'admin', 'create', 'inbound', 'InboundOrder', 13, 'RK20260204133954', NULL, '', '', '{"items": [{"remark": "", "quantity": 2, "unit_price": 0, "material_id": 30}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 13:39:55', '2026-02-04 13:39:55'),
(27, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 13, 'RK20260204133954', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 13:44:35', '2026-02-04 13:44:35'),
(28, 1, 'admin', 'create', 'inbound', 'InboundOrder', 14, 'RK20260204134503', NULL, '', '', '{"items": [{"remark": "", "quantity": 8, "unit_price": 0, "material_id": 32}, {"remark": "", "quantity": 8, "unit_price": 0, "material_id": 34}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 13:45:03', '2026-02-04 13:45:03'),
(29, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 14, 'RK20260204134503', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 13:49:54', '2026-02-04 13:49:54'),
(30, 1, 'admin', 'create', 'inbound', 'InboundOrder', 15, 'RK20260204135013', NULL, '', '', '{"items": [{"remark": "", "quantity": 4, "unit_price": 0, "material_id": 36}, {"remark": "", "quantity": 6, "unit_price": 0, "material_id": 38}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 13:50:13', '2026-02-04 13:50:13'),
(31, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 15, 'RK20260204135013', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 13:52:54', '2026-02-04 13:52:54'),
(32, 1, 'admin', 'create', 'inbound', 'InboundOrder', 16, 'RK20260204135403', NULL, '', '', '{"items": [{"remark": "", "quantity": 12, "unit_price": 0, "material_id": 42}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 13:54:03', '2026-02-04 13:54:03'),
(33, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 16, 'RK20260204135403', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 13:54:07', '2026-02-04 13:54:07'),
(34, 1, 'admin', 'create', 'inbound', 'InboundOrder', 17, 'RK20260204140457', NULL, '', '', '{"items": [{"remark": "", "quantity": 24, "unit_price": 0, "material_id": 40}], "notes": "", "remark": "", "status": "", "contact": "", "plan_id": 1, "supplier": "测试", "project_id": "10", "total_amount": 0}', '', '', 'success', '', '2026-02-04 14:04:58', '2026-02-04 14:04:58'),
(35, 1, '未知用户', 'approve', 'inbound', 'InboundOrder', 17, 'RK20260204140457', NULL, '', '', '{"comment": ""}', '', '', 'success', '', '2026-02-04 14:05:01', '2026-02-04 14:05:01');

-- Table: progress_calendars

-- Schema and data for table: progress_calendars
INSERT INTO "progress_calendars" ("id", "project_id", "name", "work_days", "holidays", "work_hours", "description", "is_default", "created_at", "updated_at") VALUES
(1, 1, '标准工作日历', '[1, 2, 3, 4, 5]', '["2024-01-01", "2024-04-04", "2024-05-01", "2024-10-01"]', '{"start": "08:00", "end": "18:00"}', '标准工作日历，周一至周五工作', true, '2025-07-01 00:00:00', '2025-07-01 00:00:00');

-- Table: progress_links

-- Schema and data for table: progress_links

-- Table: progress_projects

-- Schema and data for table: progress_projects
INSERT INTO "progress_projects" ("id", "name", "description", "start_date", "end_date", "status", "progress", "manager", "gantt_config", "network_data", "created_at", "updated_at", "created_by") VALUES
(1, '办公楼建设项目', '5层办公楼建设，包含基础工程、主体工程、装修工程', NULL, NULL, 'planning', 0, '张三', '{"tasks": [{"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149234", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 1, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u57fa\\u7840\\u5de5\\u7a0b", "type": "task", "updated_at": "2025-07-01T08:47:08.149236"}, {"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149325", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 2, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u4e3b\\u4f53\\u7ed3\\u6784", "type": "task", "updated_at": "2025-07-01T08:47:08.149325"}, {"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149365", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 3, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u88c5\\u4fee\\u5de5\\u7a0b", "type": "task", "updated_at": "2025-07-01T08:47:08.149365"}, {"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149400", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 4, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u7ae3\\u5de5\\u9a8c\\u6536", "type": "task", "updated_at": "2025-07-01T08:47:08.149400"}], "links": [], "resources": [], "calendars": [], "metadata": {"last_updated": "2025-07-01T08:47:45.747068", "version": "2.0", "format": "unified"}}', '{"tasks": [{"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149234", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 1, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u57fa\\u7840\\u5de5\\u7a0b", "type": "task", "updated_at": "2025-07-01T08:47:08.149236"}, {"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149325", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 2, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u4e3b\\u4f53\\u7ed3\\u6784", "type": "task", "updated_at": "2025-07-01T08:47:08.149325"}, {"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149365", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 3, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u88c5\\u4fee\\u5de5\\u7a0b", "type": "task", "updated_at": "2025-07-01T08:47:08.149365"}, {"actual_cost": 0, "assignee": null, "cost": 0, "created_at": "2025-07-01T08:47:08.149400", "duration": 1, "end_date": "2025-07-01 00:00:00", "id": 4, "notes": null, "parent": 0, "parent_id": null, "predecessor": null, "priority": "normal", "progress": 0, "sort_order": 0, "sortorder": 0, "start_date": "2025-06-30 00:00:00", "successor": null, "text": "\\u7ae3\\u5de5\\u9a8c\\u6536", "type": "task", "updated_at": "2025-07-01T08:47:08.149400"}], "links": [], "resources": [], "calendars": [], "metadata": {"last_updated": "2025-07-01T08:47:45.747068", "version": "2.0", "format": "unified"}}', '2025-07-01 00:00:00', '2025-07-01 00:00:00', NULL);

-- Table: progress_resources

-- Schema and data for table: progress_resources
INSERT INTO "progress_resources" ("id", "project_id", "name", "resource_type", "cost_per_day", "availability", "skills", "department", "contact_info", "notes", "created_at", "updated_at") VALUES
(1, 1, '张工程师', '人力', 500, 1, '["结构设计", "项目管理"]', '工程部', NULL, '高级工程师，负责项目技术指导', '2025-07-01 00:00:00', '2025-07-01 00:00:00'),
(2, 1, '李班长', '人力', 300, 1, '["施工管理", "质量控制"]', '施工部', NULL, '有丰富的现场施工管理经验', '2025-07-01 00:00:00', '2025-07-01 00:00:00'),
(3, 1, '挖掘机', '设备', 800, 0.8, '[]', '设备部', NULL, '日立ZX200型挖掘机', '2025-07-01 00:00:00', '2025-07-01 00:00:00');

-- Table: progress_task_resources

-- Schema and data for table: progress_task_resources

-- Table: progress_tasks

-- Schema and data for table: progress_tasks
INSERT INTO "progress_tasks" ("id", "project_id", "text", "start_date", "end_date", "duration", "progress", "task_type", "priority", "assignee", "predecessor", "successor", "notes", "cost", "actual_cost", "parent_id", "sort_order", "created_at", "updated_at") VALUES
(1, 1, '基础工程', '2025-06-30 00:00:00', '2025-07-01 00:00:00', 1, 0, 'task', 'normal', NULL, NULL, NULL, NULL, 0, 0, 0, 0, '2025-07-01 00:00:00', '2025-07-01 00:00:00'),
(2, 1, '主体结构', '2025-06-30 00:00:00', '2025-07-01 00:00:00', 1, 0, 'task', 'normal', NULL, NULL, NULL, NULL, 0, 0, 0, 0, '2025-07-01 00:00:00', '2025-07-01 00:00:00'),
(3, 1, '装修工程', '2025-06-30 00:00:00', '2025-07-01 00:00:00', 1, 0, 'task', 'normal', NULL, NULL, NULL, NULL, 0, 0, 0, 0, '2025-07-01 00:00:00', '2025-07-01 00:00:00'),
(4, 1, '竣工验收', '2025-06-30 00:00:00', '2025-07-01 00:00:00', 1, 0, 'task', 'normal', NULL, NULL, NULL, NULL, 0, 0, 0, 0, '2025-07-01 00:00:00', '2025-07-01 00:00:00');

-- Table: project_schedules

-- Schema and data for table: project_schedules
INSERT INTO "project_schedules" ("id", "project_id", "data", "created_at", "updated_at", "created_by", "updated_by") VALUES
(3, 11, '{"nodes": {}, "activities": {}}', '2026-01-30 04:33:24', '2026-01-30 04:33:24', NULL, NULL);

-- Table: projects

-- Schema and data for table: projects
INSERT INTO "projects" ("id", "name", "code", "location", "start_date", "end_date", "description", "manager", "contact", "budget", "status", "parent_id", "level", "path", "progress_percentage") VALUES
(11, '腈纶厂长丝项目', 'QLCCSXM', '腈纶厂', '2025-01-30 00:00:00', '2025-08-30 00:00:00', '', '', '', '', 'completed', NULL, 0, '/11/', '0.00'),
(9, '腈纶厂长丝项目空调机组', 'PJ-20250622009', '腈纶厂', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/9/', '0.00'),
(10, '线性装置聚烯烃弹性体生产优化改造', 'PJ-20260121001', '浦江', '2026-01-20 00:00:00', '2026-05-30 00:00:00', '？\n\n', '鞠磊', '', '0', 'active', NULL, 0, '/10/', '0.78'),
(8, '腈纶厂长丝项目电话系统', 'PJ-20250622008', '腈纶厂', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/8/', '0.00'),
(7, '腈纶厂长丝项目门禁系统', 'PJ-20250622007', '腈纶厂', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/7/', '0.00'),
(6, '腈纶厂长丝项目广播系统', 'PJ-20250622006', '腈纶厂', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/6/', '0.00'),
(5, '腈纶厂长丝项目视频监控', 'PJ-20250622005', '腈纶厂', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/5/', '0.00'),
(4, '腈纶厂长丝项目火灾', 'PJ-20250622004', '腈纶厂', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/4/', '0.00'),
(2, '腈纶厂长丝项目仪表专有设备区', 'PJ-20250622002', '腈纶', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/2/', '0.00'),
(1, '腈纶厂长丝项目仪表附属设备区', 'PJ-20250622001', '腈纶', '2025-06-30 00:00:00', '2025-09-30 00:00:00', '', '鞠磊', '18045969960', '0', 'completed', 11, 1, '/11/1/', '0.00');

-- Table: requisition_items

-- Schema and data for table: requisition_items
INSERT INTO "requisition_items" ("id", "requisition_id", "stock_id", "material_id", "requested_quantity", "approved_quantity", "actual_quantity", "status", "remark", "created_at", "updated_at") VALUES
(1, 1, 2, 4, '16.000', '16.000', '16.000', 'issued', '', '2026-02-03 09:01:27', '2026-02-03 09:01:33'),
(2, 1, 1, 2, '20.000', '20.000', '20.000', 'issued', '', '2026-02-03 09:01:27', '2026-02-03 09:01:33'),
(3, 2, 1, 2, '5.000', '5.000', '5.000', 'issued', '', '2026-02-04 05:42:43', '2026-02-04 06:36:01');

-- Table: requisitions

-- Schema and data for table: requisitions
INSERT INTO "requisitions" ("id", "requisition_no", "project_name", "applicant", "department", "status", "created_at", "remark", "approved_at", "approved_by", "urgent", "purpose", "issued_by", "issued_at", "project_id", "updated_at", "plan_id") VALUES
(1, 'CK20260203001', NULL, 'admin', '', 'issued', '2026-02-03 09:01:27', '', '2026-02-03 09:01:30', 'admin', 0, '', 'admin', '2026-02-03 09:01:33', 10, '2026-02-03 09:01:27', NULL),
(2, 'CK20260204001', NULL, 'admin', '', 'issued', '2026-02-04 05:42:43', '', '2026-02-04 05:42:50', 'admin', 0, '测试审核', 'admin', '2026-02-04 06:36:01', 10, '2026-02-04 05:42:43', NULL);

-- Table: resources

-- Schema and data for table: resources
INSERT INTO "resources" ("id", "project_id", "name", "type", "unit", "quantity", "cost_per_unit", "color", "is_active", "created_at", "updated_at") VALUES
(1, 1, '高级工程师', 'labor', '人/d', '10', '800', '#409EFF', true, '2026-01-30 13:15:13', '2026-01-30 13:16:05'),
(2, 10, '技工', 'labor', '人/d', '2', '0', '#409eff', true, '2026-01-31 02:27:46', '2026-01-31 02:27:46'),
(3, 10, '力工', 'labor', '人/d', '5', '150', '#409eff', true, '2026-01-31 02:28:21', '2026-01-31 02:28:21'),
(4, 10, '電焊機', 'equipment', '人/d', '5', '100', 'rgb(255, 220, 64)', true, '2026-01-31 13:07:52', '2026-01-31 13:07:52');

-- Table: roles

-- Schema and data for table: roles
INSERT INTO "roles" ("id", "name", "description", "permissions", "created_at") VALUES
(2, '保管员', '保管员角色', 'material_view,stocklog_view,stock_view,stock_out,stock_in,stock_export,stock_edit,requisition_view,requisition_issue,requisition_export,inbound_view,inbound_export,inbound_approve,system_statistics,system_report,system_log,system_activities', NULL),
(6, '材料员', '材料员', 'stock,stock_view,stock_in,stock_out,stock_edit,stock_delete,stocklog_view,stock_export,requisition,requisition_view,requisition_create,requisition_edit,requisition_delete,requisition_approve,requisition_issue,requisition_export,inbound,inbound_view,inbound_create,inbound_edit,inbound_delete,inbound_approve,inbound_export', '0001-01-01 00:00:00'),
(1, 'admin', '管理员角色', 'admin,audit_view', NULL),
(4, '分包材料员', '分包材料员', 'project_view,material_view,stocklog_view,stock_view,stock_export,requisition_view,requisition_delete,requisition_create', '0001-01-01 00:00:00'),
(3, '施工员', '施工员', 'material,material_view,material_create,material_edit,material_delete,material_import,material_export,stock,stock_view,stock_in,stock_out,stock_edit,stock_delete,stocklog_view,stock_export,requisition,requisition_view,requisition_create,requisition_edit,requisition_delete,requisition_approve,requisition_issue,requisition_export,inbound,inbound_view,inbound_create,inbound_edit,inbound_delete,inbound_approve,inbound_export,project,project_view,project_create,project_edit,project_delete,project_export,construction_log,construction_log_view,construction_log_create,construction_log_edit,construction_log_delete,construction_log_export,progress,progress_view,progress_create,progress_edit,progress_delete,progress_export,user_view,role_view,system', '0001-01-01 00:00:00'),
(5, '项目经理', '项目经理', 'material,material_view,material_create,material_edit,material_delete,material_import,material_export,stock,stock_view,stock_in,stock_out,stock_edit,stock_delete,stocklog_view,stock_export,requisition,requisition_view,requisition_create,requisition_edit,requisition_delete,requisition_approve,requisition_issue,requisition_export,inbound,inbound_view,inbound_create,inbound_edit,inbound_delete,inbound_approve,inbound_export,project,project_view,project_create,project_edit,project_delete,project_export,construction_log,construction_log_view,construction_log_create,construction_log_edit,construction_log_delete,construction_log_export,progress,progress_view,progress_create,progress_edit,progress_delete,progress_export,user_view,role_view,log_view,system', '0001-01-01 00:00:00');

-- Table: stock_logs

-- Schema and data for table: stock_logs
INSERT INTO "stock_logs" ("id", "stock_id", "type", "quantity", "quantity_before", "quantity_after", "source_type", "source_id", "source_no", "project_id", "material_id", "user_id", "remark", "created_at") VALUES
(1, 1, 'in', '40.000', '0.000', '40.000', 'inbound', 1, 'RK20260203085209', 10, 2, 1, '入库 40.00 个，备注：入库单 RK20260203085209', '2026-02-03 08:52:13'),
(2, 2, 'in', '16.000', '0.000', '16.000', 'inbound', 1, 'RK20260203085209', 10, 4, 1, '入库 16.00 个，备注：入库单 RK20260203085209', '2026-02-03 08:52:13'),
(3, 2, 'out', '16.000', '16.000', '0.000', 'requisition', 1, 'CK20260203001', 10, 4, 1, '出库单发放：CK20260203001，出库 16.00', '2026-02-03 09:01:33'),
(4, 1, 'out', '20.000', '40.000', '20.000', 'requisition', 1, 'CK20260203001', 10, 2, 1, '出库单发放：CK20260203001，出库 20.00', '2026-02-03 09:01:33'),
(5, 3, 'in', '260.000', '0.000', '260.000', 'inbound', 3, 'RK20260204051229', 10, 6, 1, '入库 260.00 米，备注：入库单 RK20260204051229', '2026-02-04 06:30:59'),
(6, 4, 'in', '2.000', '0.000', '2.000', 'inbound', 3, 'RK20260204051229', 10, 8, 1, '入库 2.00 个，备注：入库单 RK20260204051229', '2026-02-04 06:30:59'),
(7, 1, 'in', '40.000', '20.000', '60.000', 'inbound', 2, 'RK20260204043801', 10, 2, 1, '入库 40.00 个，备注：入库单 RK20260204043801', '2026-02-04 06:31:18'),
(8, 2, 'in', '16.000', '0.000', '16.000', 'inbound', 2, 'RK20260204043801', 10, 4, 1, '入库 16.00 个，备注：入库单 RK20260204043801', '2026-02-04 06:31:18'),
(9, 1, 'out', '5.000', '60.000', '55.000', 'requisition', 2, 'CK20260204001', 10, 2, 1, '出库单发放：CK20260204001，出库 5.00', '2026-02-04 06:36:01'),
(10, 5, 'in', '6.000', '0.000', '6.000', 'inbound', 4, 'RK20260204063628', 10, 10, 1, '入库 6.00 个，备注：入库单 RK20260204063628', '2026-02-04 06:36:32'),
(11, 6, 'in', '12.000', '0.000', '12.000', 'inbound', 4, 'RK20260204063628', 10, 12, 1, '入库 12.00 米，备注：入库单 RK20260204063628', '2026-02-04 06:36:32'),
(12, 7, 'in', '198.000', '0.000', '198.000', 'inbound', 5, 'RK20260204064523', 10, 14, 1, '入库 198.00 米，备注：入库单 RK20260204064523', '2026-02-04 06:45:26'),
(13, 8, 'in', '4.000', '0.000', '4.000', 'inbound', 5, 'RK20260204064523', 10, 16, 1, '入库 4.00 个，备注：入库单 RK20260204064523', '2026-02-04 06:45:26'),
(14, 7, 'in', '2.000', '198.000', '200.000', 'inbound', 6, 'RK20260204064552', 10, 14, 1, '入库 2.00 米，备注：入库单 RK20260204064552', '2026-02-04 06:45:56'),
(15, 9, 'in', '2.000', '0.000', '2.000', 'inbound', 8, 'RK20260204075456', 10, 18, 1, '入库 2.00 个，备注：入库单 RK20260204075456', '2026-02-04 07:55:00'),
(16, 10, 'in', '78.000', '0.000', '78.000', 'inbound', 8, 'RK20260204075456', 10, 20, 1, '入库 78.00 米，备注：入库单 RK20260204075456', '2026-02-04 07:55:00'),
(17, 11, 'in', '13.000', '0.000', '13.000', 'inbound', 9, 'RK20260204085325', 10, 22, 1, '入库 13.00 个，备注：入库单 RK20260204085325', '2026-02-04 08:53:32'),
(18, 12, 'in', '6.000', '0.000', '6.000', 'inbound', 10, 'RK20260204085833', 10, 24, 1, '入库 6.00 个，备注：入库单 RK20260204085833', '2026-02-04 08:58:41'),
(19, 13, 'in', '560.000', '0.000', '560.000', 'inbound', 11, 'RK20260204133723', 10, 26, 1, '入库 560.00 米，备注：入库单 RK20260204133723', '2026-02-04 13:37:32'),
(20, 14, 'in', '100.000', '0.000', '100.000', 'inbound', 12, 'RK20260204133807', 10, 28, 1, '入库 100.00 米，备注：入库单 RK20260204133807', '2026-02-04 13:39:09'),
(21, 15, 'in', '2.000', '0.000', '2.000', 'inbound', 13, 'RK20260204133954', 10, 30, 1, '入库 2.00 个，备注：入库单 RK20260204133954', '2026-02-04 13:44:35'),
(22, 16, 'in', '8.000', '0.000', '8.000', 'inbound', 14, 'RK20260204134503', 10, 32, 1, '入库 8.00 条，备注：入库单 RK20260204134503', '2026-02-04 13:49:54'),
(23, 17, 'in', '8.000', '0.000', '8.000', 'inbound', 14, 'RK20260204134503', 10, 34, 1, '入库 8.00 条，备注：入库单 RK20260204134503', '2026-02-04 13:49:54'),
(24, 18, 'in', '4.000', '0.000', '4.000', 'inbound', 15, 'RK20260204135013', 10, 36, 1, '入库 4.00 条，备注：入库单 RK20260204135013', '2026-02-04 13:52:54'),
(25, 19, 'in', '6.000', '0.000', '6.000', 'inbound', 15, 'RK20260204135013', 10, 38, 1, '入库 6.00 个，备注：入库单 RK20260204135013', '2026-02-04 13:52:54');

-- Table: stock_op_logs

-- Schema and data for table: stock_op_logs
INSERT INTO "stock_op_logs" ("id", "user_id", "op_type", "stock_id", "log_id", "detail", "time", "created_at", "updated_at") VALUES
(80, 1, 'in', 1, NULL, '入库 40.00 个，备注：入库单 RK20260203085209', '2026-02-03 08:52:13', '2026-02-03 08:52:13', '2026-02-03 08:52:13'),
(81, 1, 'in', 2, NULL, '入库 16.00 个，备注：入库单 RK20260203085209', '2026-02-03 08:52:13', '2026-02-03 08:52:13', '2026-02-03 08:52:13'),
(82, 1, 'out', 2, 1, '出库单发放：CK20260203001，出库 16.00', '2026-02-03 09:01:33', '2026-02-03 09:01:33', '2026-02-03 09:01:33'),
(83, 1, 'out', 1, 1, '出库单发放：CK20260203001，出库 20.00', '2026-02-03 09:01:33', '2026-02-03 09:01:33', '2026-02-03 09:01:33'),
(84, 1, 'in', 3, NULL, '入库 260.00 米，备注：入库单 RK20260204051229', '2026-02-04 06:30:59', '2026-02-04 06:30:59', '2026-02-04 06:30:59'),
(85, 1, 'in', 4, NULL, '入库 2.00 个，备注：入库单 RK20260204051229', '2026-02-04 06:30:59', '2026-02-04 06:30:59', '2026-02-04 06:30:59'),
(86, 1, 'in', 1, NULL, '入库 40.00 个，备注：入库单 RK20260204043801', '2026-02-04 06:31:18', '2026-02-04 06:31:18', '2026-02-04 06:31:18'),
(87, 1, 'in', 2, NULL, '入库 16.00 个，备注：入库单 RK20260204043801', '2026-02-04 06:31:18', '2026-02-04 06:31:18', '2026-02-04 06:31:18'),
(88, 1, 'out', 1, 2, '出库单发放：CK20260204001，出库 5.00', '2026-02-04 06:36:01', '2026-02-04 06:36:01', '2026-02-04 06:36:01'),
(89, 1, 'in', 5, NULL, '入库 6.00 个，备注：入库单 RK20260204063628', '2026-02-04 06:36:32', '2026-02-04 06:36:32', '2026-02-04 06:36:32'),
(90, 1, 'in', 6, NULL, '入库 12.00 米，备注：入库单 RK20260204063628', '2026-02-04 06:36:32', '2026-02-04 06:36:32', '2026-02-04 06:36:32'),
(91, 1, 'in', 7, NULL, '入库 198.00 米，备注：入库单 RK20260204064523', '2026-02-04 06:45:26', '2026-02-04 06:45:26', '2026-02-04 06:45:26'),
(92, 1, 'in', 8, NULL, '入库 4.00 个，备注：入库单 RK20260204064523', '2026-02-04 06:45:26', '2026-02-04 06:45:26', '2026-02-04 06:45:26'),
(93, 1, 'in', 7, NULL, '入库 2.00 米，备注：入库单 RK20260204064552', '2026-02-04 06:45:56', '2026-02-04 06:45:56', '2026-02-04 06:45:56'),
(94, 1, 'in', 9, NULL, '入库 2.00 个，备注：入库单 RK20260204075456', '2026-02-04 07:55:00', '2026-02-04 07:55:00', '2026-02-04 07:55:00'),
(95, 1, 'in', 10, NULL, '入库 78.00 米，备注：入库单 RK20260204075456', '2026-02-04 07:55:00', '2026-02-04 07:55:00', '2026-02-04 07:55:00'),
(96, 1, 'in', 11, NULL, '入库 13.00 个，备注：入库单 RK20260204085325', '2026-02-04 08:53:32', '2026-02-04 08:53:32', '2026-02-04 08:53:32'),
(97, 1, 'in', 12, NULL, '入库 6.00 个，备注：入库单 RK20260204085833', '2026-02-04 08:58:41', '2026-02-04 08:58:41', '2026-02-04 08:58:41'),
(98, 1, 'in', 13, NULL, '入库 560.00 米，备注：入库单 RK20260204133723', '2026-02-04 13:37:32', '2026-02-04 13:37:32', '2026-02-04 13:37:32'),
(99, 1, 'in', 14, NULL, '入库 100.00 米，备注：入库单 RK20260204133807', '2026-02-04 13:39:09', '2026-02-04 13:39:09', '2026-02-04 13:39:09'),
(100, 1, 'in', 15, NULL, '入库 2.00 个，备注：入库单 RK20260204133954', '2026-02-04 13:44:35', '2026-02-04 13:44:35', '2026-02-04 13:44:35'),
(101, 1, 'in', 16, NULL, '入库 8.00 条，备注：入库单 RK20260204134503', '2026-02-04 13:49:54', '2026-02-04 13:49:54', '2026-02-04 13:49:54'),
(102, 1, 'in', 17, NULL, '入库 8.00 条，备注：入库单 RK20260204134503', '2026-02-04 13:49:54', '2026-02-04 13:49:54', '2026-02-04 13:49:54'),
(103, 1, 'in', 18, NULL, '入库 4.00 条，备注：入库单 RK20260204135013', '2026-02-04 13:52:54', '2026-02-04 13:52:54', '2026-02-04 13:52:54'),
(104, 1, 'in', 19, NULL, '入库 6.00 个，备注：入库单 RK20260204135013', '2026-02-04 13:52:54', '2026-02-04 13:52:54', '2026-02-04 13:52:54');

-- Table: stocks

-- Schema and data for table: stocks
INSERT INTO "stocks" ("id", "project_id", "material_id", "warehouse_id", "quantity", "safety_stock", "location", "unit_cost", "updated_at", "created_at") VALUES
(3, 10, 6, NULL, '260.000', '0.000', '', '0.00', '2026-02-04 06:30:59', '2026-02-04 06:30:59'),
(4, 10, 8, NULL, '2.000', '0.000', '', '0.00', '2026-02-04 06:30:59', '2026-02-04 06:30:59'),
(2, 10, 4, NULL, '16.000', '0.000', '', '0.00', '2026-02-03 09:01:33', '2026-02-03 08:52:13'),
(1, 10, 2, NULL, '55.000', '0.000', '', '0.00', '2026-02-04 06:36:01', '2026-02-03 08:52:13'),
(5, 10, 10, NULL, '6.000', '0.000', '', '0.00', '2026-02-04 06:36:32', '2026-02-04 06:36:32'),
(6, 10, 12, NULL, '12.000', '0.000', '', '0.00', '2026-02-04 06:36:32', '2026-02-04 06:36:32'),
(8, 10, 16, NULL, '4.000', '0.000', '', '0.00', '2026-02-04 06:45:26', '2026-02-04 06:45:26'),
(7, 10, 14, NULL, '200.000', '0.000', '', '0.00', '2026-02-04 06:45:26', '2026-02-04 06:45:26'),
(9, 10, 18, NULL, '2.000', '0.000', '', '0.00', '2026-02-04 07:55:00', '2026-02-04 07:55:00'),
(10, 10, 20, NULL, '78.000', '0.000', '', '0.00', '2026-02-04 07:55:00', '2026-02-04 07:55:00'),
(11, 10, 22, NULL, '13.000', '0.000', '', '0.00', '2026-02-04 08:53:31', '2026-02-04 08:53:31'),
(12, 10, 24, NULL, '6.000', '0.000', '', '0.00', '2026-02-04 08:58:41', '2026-02-04 08:58:41'),
(13, 10, 26, NULL, '560.000', '0.000', '', '0.00', '2026-02-04 13:37:32', '2026-02-04 13:37:32'),
(14, 10, 28, NULL, '100.000', '0.000', '', '0.00', '2026-02-04 13:39:09', '2026-02-04 13:39:09'),
(15, 10, 30, NULL, '2.000', '0.000', '', '0.00', '2026-02-04 13:44:35', '2026-02-04 13:44:35'),
(16, 10, 32, NULL, '8.000', '0.000', '', '0.00', '2026-02-04 13:49:54', '2026-02-04 13:49:54'),
(17, 10, 34, NULL, '8.000', '0.000', '', '0.00', '2026-02-04 13:49:54', '2026-02-04 13:49:54'),
(18, 10, 36, NULL, '4.000', '0.000', '', '0.00', '2026-02-04 13:52:54', '2026-02-04 13:52:54'),
(19, 10, 38, NULL, '6.000', '0.000', '', '0.00', '2026-02-04 13:52:54', '2026-02-04 13:52:54'),
(20, 10, 42, NULL, '12.000', '0.000', '', '0.00', '2026-02-04 13:54:07', '2026-02-04 13:54:07'),
(21, 10, 40, NULL, '24.000', '0.000', '', '0.00', '2026-02-04 14:05:01', '2026-02-04 14:05:01');

-- Table: system_activity

-- Schema and data for table: system_activity

-- Table: system_backup

-- Schema and data for table: system_backup
INSERT INTO "system_backup" ("id", "filename", "filepath", "size", "status", "created_by", "created_at", "description") VALUES
(1, 'backup_20250624_081256.sql', '/app/backup_20250624_081256.sql', 147456, 'completed', 'admin', '2025-06-24 08:12:56', '系统自动备份 - 20250624_081256'),
(2, 'backup_20250626_055730.sql', '/app/backup_20250626_055730.sql', 147456, 'completed', 'admin', '2025-06-26 05:57:30', '系统自动备份 - 20250626_055730'),
(3, 'backup_20250721_003210.sql', '/app/backup_20250721_003210.sql', 442368, 'completed', 'admin', '2025-07-21 00:32:10', '系统自动备份 - 20250721_003210'),
(4, 'backup_20251216_070717.sql', '/app/backup_20251216_070717.sql', 540672, 'completed', 'admin', '2025-12-16 07:07:17', '系统自动备份 - 20251216_070717'),
(10, 'backup_20260124_145429.sql', 'backup_20260124_145429.sql', 408054, 'completed', 'julei', '2026-01-24 14:54:29', 'PostgreSQL数据库备份 - 20260124_145429'),
(11, 'backup_20260126_072844.sql', 'backup_20260126_072844.sql', 416615, 'completed', 'julei', '2026-01-26 07:28:44', 'PostgreSQL数据库备份 - 20260126_072844'),
(12, 'backup_20260127_134856.sql', 'backup_20260127_134856.sql', 431410, 'completed', 'admin', '2026-01-27 13:48:56', 'PostgreSQL数据库备份 - 20260127_134856'),
(13, 'backup_20260128_043930.sql', 'backup_20260128_043930.sql', 434498, 'completed', 'julei', '2026-01-28 04:39:30', 'PostgreSQL数据库备份 - 20260128_043930'),
(14, 'backup_20260130_071909.sql', 'backup_20260130_071909.sql', 460837, 'completed', 'admin', '2026-01-30 07:19:09', 'PostgreSQL数据库备份 - 20260130_071909');

-- Table: system_config

-- Schema and data for table: system_config
INSERT INTO "system_config" ("id", "key", "value", "type", "description", "created_at", "updated_at") VALUES
(1, 'system_name', '项目管理系统', 'string', '系统配置: system_name', '2025-06-21 00:00:00', '2026-02-03 00:00:00'),
(8, 'system_short_name', 'MMS', 'string', '', NULL, '2026-02-03 00:00:00'),
(10, 'token_expiry', '72', 'integer', '', NULL, '2026-02-03 00:00:00'),
(4, 'max_file_size', '5', 'integer', '最大文件上传大小(MB)', NULL, '2026-02-03 00:00:00'),
(6, 'max_upload_count', '10', 'integer', '单次最多上传文件数量', NULL, '2026-02-03 00:00:00'),
(5, 'allowed_file_types', 'jpg,jpeg,png,gif,bmp,webp,svg', 'string', '允许上传的文件类型', NULL, '2026-02-03 00:00:00'),
(9, 'password_min_length', '6', 'integer', '', NULL, '2026-02-03 00:00:00'),
(7, 'enable_captcha', 'false', 'boolean', '', NULL, '2026-02-03 00:00:00'),
(3, 'upload_directory', 'static/uploads', 'string', '文件上传目录路径', NULL, '2026-02-03 00:00:00'),
(2, 'session_timeout', '30', 'integer', '系统配置: session_timeout', '2025-06-21 00:00:00', '2026-02-03 00:00:00');

-- Table: system_logs

-- Schema and data for table: system_logs
,
INSERT INTO "system_logs" ("id", "level", "message", "module", "user_id", "ip_address", "created_at") VALUES
(530, 'INFO', '执行数据库备份: backup_20260126_072844.sql', 'backup', 3, '127.0.0.1', NULL),
(531, 'INFO', '执行数据库备份: backup_20260127_134856.sql', 'backup', 1, '192.168.18.1', NULL),
(532, 'INFO', '执行数据库备份: backup_20260128_043930.sql', 'backup', 3, '192.168.18.1', NULL),
(533, 'INFO', '执行数据库备份: backup_20260130_071909.sql', 'backup', 1, '192.168.18.1', NULL),
(1, 'WARNING', '清空系统日志', 'system', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(2, 'INFO', '执行数据库备份: backup_20250624_081150.sql', 'backup', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(3, 'INFO', '执行数据库备份: backup_20250624_081256.sql', 'backup', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(4, 'INFO', '更新系统设置: [''system_name'', ''session_timeout'']', 'system', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(5, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(6, 'INFO', '创建入库单: RK20250624083008', 'inbound', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(7, 'WARNING', '删除入库单: RK20250624083008', 'inbound', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(8, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(9, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-24 00:00:00'),
(10, 'INFO', '创建入库单: RK20250625030939', 'inbound', 1, '172.18.0.1', '2025-06-25 00:00:00'),
(11, 'INFO', '审批入库单: RK20250625030939，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-06-25 00:00:00'),
(12, 'INFO', '创建入库单: RK20250625031154', 'inbound', 1, '172.18.0.1', '2025-06-25 00:00:00'),
(13, 'INFO', '审批入库单: RK20250625031154，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-06-25 00:00:00'),
(14, 'INFO', '库存调整: 镀锌钢管 (类型:set, 新库存:360.0)', 'stock', 1, '172.18.0.1', '2025-06-25 00:00:00'),
(15, 'INFO', '执行数据库备份: backup_20250626_055730.sql', 'backup', 1, '172.18.0.1', '2025-06-26 00:00:00'),
(16, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(17, 'INFO', '编辑项目: 腈纶厂长丝项目火灾 (编号:PJ-20250622004)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(18, 'INFO', '编辑项目: 腈纶厂长丝项目视频监控 (编号:PJ-20250622005)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(19, 'INFO', '编辑项目: 腈纶厂长丝项目广播系统 (编号:PJ-20250622006)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(20, 'INFO', '编辑项目: 腈纶厂长丝项目门禁系统 (编号:PJ-20250622007)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(21, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(22, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(23, 'WARNING', '重置用户密码: julei', 'auth', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(24, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(25, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-27 00:00:00'),
(26, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(27, 'INFO', '编辑项目: 腈纶厂长丝项目仪表附属设备区 (编号:PJ-20250622001)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(28, 'INFO', '编辑项目: 腈纶厂长丝项目仪表专有设备区 (编号:PJ-20250622002)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(29, 'INFO', '编辑项目: 腈纶厂长丝项目电信 (编号:PJ-20250622003)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(30, 'INFO', '编辑项目: 腈纶厂长丝项目火灾 (编号:PJ-20250622004)', 'project', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(31, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-27 00:00:00'),
(32, 'INFO', '新增项目: 第三方斯蒂芬 (编号:PJ-20250628001)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(33, 'INFO', '新增项目: 特色他 (编号:PJ-20250628002)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(34, 'WARNING', '删除项目: 第三方斯蒂芬 (编号:PJ-20250628001)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(35, 'WARNING', '删除项目: 特色他 (编号:PJ-20250628002)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(36, 'INFO', '新增项目: 特色他 (编号:PJ-20250628001)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(37, 'WARNING', '删除项目: 特色他 (编号:PJ-20250628001)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(38, 'INFO', '新增项目: 特色他 (编号:PJ-20250628001)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(39, 'INFO', '新增项目: 地方的司法 (编号:PJ-20250628002)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(40, 'WARNING', '删除项目: 特色他 (编号:PJ-20250628001)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(41, 'WARNING', '删除项目: 地方的司法 (编号:PJ-20250628002)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(42, 'INFO', '新增角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(43, 'INFO', '编辑用户: julei', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(44, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(45, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(46, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(47, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(48, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(49, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(50, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(51, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(52, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(53, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(54, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(55, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(56, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(57, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(58, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(59, 'INFO', '编辑角色: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(60, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(61, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(62, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(63, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(64, 'INFO', '新增项目: rerere (编号:PJ-20250628001)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(65, 'INFO', '编辑项目: rerere (编号:PJ-20250628001)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(66, 'WARNING', '删除项目: rerere (编号:PJ-20250628001)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(67, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(68, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(69, 'INFO', '新增角色: 材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(70, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(71, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(72, 'INFO', '编辑角色: 施工员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(73, 'INFO', '编辑角色: 保管员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(74, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(75, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(76, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(77, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(78, 'INFO', '编辑项目: 腈纶厂长丝项目火灾 (编号:PJ-20250622004)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(79, 'INFO', '编辑项目: 腈纶厂长丝项目电信 (编号:PJ-20250622003)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(80, 'INFO', '编辑项目: 腈纶厂长丝项目仪表专有设备区 (编号:PJ-20250622002)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(81, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(82, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(83, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(84, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(85, 'INFO', '编辑项目: 腈纶厂长丝项目视频监控 (编号:PJ-20250622005)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(86, 'INFO', '编辑项目: 腈纶厂长丝项目广播系统 (编号:PJ-20250622006)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(87, 'INFO', '编辑项目: 腈纶厂长丝项目门禁系统 (编号:PJ-20250622007)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(88, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(89, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(90, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(91, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(92, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(93, 'INFO', '创建入库单: RK20250628080103', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(94, 'INFO', '创建入库单: RK20250628080309', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(95, 'INFO', '创建入库单: RK20250628080541', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(96, 'INFO', '创建入库单: RK20250628080726', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00');
,
INSERT INTO "system_logs" ("id", "level", "message", "module", "user_id", "ip_address", "created_at") VALUES
(97, 'INFO', '创建入库单: RK20250628080840', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(98, 'INFO', '创建入库单: RK20250628080935', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(99, 'INFO', '创建入库单: RK20250628081048', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(100, 'INFO', '创建入库单: RK20250628081106', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(101, 'INFO', '审批入库单: RK20250628080103，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(102, 'INFO', '审批入库单: RK20250628080309，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(103, 'INFO', '审批入库单: RK20250628080541，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(104, 'INFO', '审批入库单: RK20250628080726，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(105, 'INFO', '审批入库单: RK20250628080840，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(106, 'INFO', '审批入库单: RK20250628080935，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(107, 'INFO', '审批入库单: RK20250628081048，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(108, 'INFO', '审批入库单: RK20250628081106，结果: 批准', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(109, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(110, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(111, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(112, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(113, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(114, 'INFO', '创建入库单: RK20250628085141', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(115, 'INFO', '创建入库单: RK20250628090514', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(116, 'INFO', '创建入库单: RK20250628090715', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(117, 'INFO', '编辑入库单: RK20250628090715', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(118, 'INFO', '编辑入库单: RK20250628090715', 'inbound', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(119, 'INFO', '新增申请单: CK20250628001 (项目:腈纶厂长丝项目仪表附属设备区)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(120, 'INFO', '新增申请单: CK20250628002 (项目:腈纶厂长丝项目仪表专有设备区)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(121, 'INFO', '新增申请单: CK20250628003 (项目:腈纶厂长丝项目空调机组)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(122, 'INFO', '审核申请单: CK20250628003 (项目:腈纶厂长丝项目空调机组)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(123, 'INFO', '审核申请单: CK20250628002 (项目:腈纶厂长丝项目仪表专有设备区)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(124, 'INFO', '审核申请单: CK20250628001 (项目:腈纶厂长丝项目仪表附属设备区)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(125, 'INFO', '发放申请单: CK20250628003 (项目:腈纶厂长丝项目空调机组)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(126, 'INFO', '发放申请单: CK20250628001 (项目:腈纶厂长丝项目仪表附属设备区)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(127, 'INFO', '新增申请单: CK20250628004 (项目:腈纶厂长丝项目仪表专有设备区)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(128, 'INFO', '审核申请单: CK20250628004 (项目:腈纶厂长丝项目仪表专有设备区)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(129, 'INFO', '发放申请单: CK20250628004 (项目:腈纶厂长丝项目仪表专有设备区)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(130, 'WARNING', '删除申请单: CK20250628002 (项目:腈纶厂长丝项目仪表专有设备区)', 'requisition', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(131, 'INFO', '新增用户: libo', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(132, 'INFO', '编辑用户: libo', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(133, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(134, 'INFO', '用户登录: libo', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(135, 'INFO', '用户登出', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(136, 'INFO', '用户登录: libo', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(137, 'INFO', '用户登出', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(138, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(139, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(140, 'INFO', '用户登录: libo', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(141, 'INFO', '用户登出', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(142, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(143, 'INFO', '编辑用户: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(144, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(145, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(146, 'INFO', '编辑角色: 材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(147, 'INFO', '编辑角色: 保管员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(148, 'INFO', '编辑角色: 施工员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(149, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(150, 'INFO', '新增用户: gjw', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(151, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(152, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(153, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(154, 'INFO', '用户登录: libo', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(155, 'INFO', '用户登出', 'auth', 4, '172.18.0.1', '2025-06-28 00:00:00'),
(156, 'INFO', '用户登录: wqs', 'auth', 2, '172.18.0.1', '2025-06-28 00:00:00'),
(157, 'INFO', '用户登出', 'auth', 2, '172.18.0.1', '2025-06-28 00:00:00'),
(158, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(159, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(160, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(161, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(162, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(163, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(164, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(165, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(166, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(167, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(168, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(169, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(170, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(171, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(172, 'INFO', '编辑用户: gjw', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(173, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(174, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(175, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(176, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(177, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(178, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(179, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-28 00:00:00'),
(180, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(181, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(182, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(183, 'INFO', '编辑项目: 腈纶厂长丝项目仪表附属设备区 (编号:PJ-20250622001)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(184, 'INFO', '编辑项目: 腈纶厂长丝项目仪表专有设备区 (编号:PJ-20250622002)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(185, 'INFO', '编辑项目: 腈纶厂长丝项目电信 (编号:PJ-20250622003)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(186, 'INFO', '编辑项目: 腈纶厂长丝项目火灾 (编号:PJ-20250622004)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(187, 'INFO', '编辑项目: 腈纶厂长丝项目视频监控 (编号:PJ-20250622005)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(275, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(188, 'INFO', '编辑项目: 腈纶厂长丝项目广播系统 (编号:PJ-20250622006)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(189, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(190, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(191, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-28 00:00:00'),
(192, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(193, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(194, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(195, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00');
,
INSERT INTO "system_logs" ("id", "level", "message", "module", "user_id", "ip_address", "created_at") VALUES
(196, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-28 00:00:00'),
(197, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-29 00:00:00'),
(198, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(199, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(200, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(201, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(202, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-29 00:00:00'),
(203, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(204, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(205, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(206, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(207, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(208, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(209, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(210, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(211, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-29 00:00:00'),
(212, 'INFO', '用户登录: libo', 'auth', 4, '172.18.0.1', '2025-06-29 00:00:00'),
(213, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(214, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(215, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(216, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(217, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(218, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(219, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(220, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(221, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(222, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(223, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(224, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(225, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(226, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(227, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-29 00:00:00'),
(228, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(229, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(230, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(231, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(232, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-29 00:00:00'),
(233, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-29 00:00:00'),
(234, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(235, 'INFO', '新增申请单: CK20250629002 (项目ID:None)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(236, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(237, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(238, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(239, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(240, 'INFO', '新增申请单: CK20250629003 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(241, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(242, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(243, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(244, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(245, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(246, 'INFO', '新增申请单: CK20250629004 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(247, 'WARNING', '删除申请单: CK20250629004 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(248, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(249, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(250, 'WARNING', '删除申请单: CK20250629002 (项目ID:None)', 'requisition', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(251, 'WARNING', '删除申请单: CK20250629001 (项目ID:None)', 'requisition', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(252, 'WARNING', '删除申请单: CK20250629003 (项目ID:1)', 'requisition', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(253, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-29 00:00:00'),
(254, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(255, 'INFO', '新增申请单: CK20250629001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(256, 'WARNING', '删除申请单: CK20250629001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(257, 'INFO', '新增申请单: CK20250629001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(258, 'WARNING', '删除申请单: CK20250629001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-29 00:00:00'),
(259, 'INFO', '审批入库单: RK20250628090715，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(260, 'INFO', '审批入库单: RK20250628090514，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(261, 'INFO', '审批入库单: RK20250628085141，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(262, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(263, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-06-30 00:00:00'),
(264, 'INFO', '新增申请单: CK20250630001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-06-30 00:00:00'),
(265, 'INFO', '新增申请单: CK20250630002 (项目ID:2)', 'requisition', 5, '172.18.0.1', '2025-06-30 00:00:00'),
(266, 'INFO', '新增申请单: CK20250630003 (项目ID:9)', 'requisition', 5, '172.18.0.1', '2025-06-30 00:00:00'),
(267, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-06-30 00:00:00'),
(268, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(269, 'INFO', '审核申请单: CK20250630003 (项目ID:9)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(270, 'INFO', '审核申请单: CK20250630002 (项目ID:2)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(271, 'INFO', '审核申请单: CK20250630001 (项目ID:1)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(272, 'INFO', '发放申请单: CK20250630003 (项目ID:9)', 'requisition', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(273, 'INFO', '发放申请单: CK20250630002 (项目ID:2)', 'requisition', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(274, 'INFO', '发放申请单: CK20250630001 (项目ID:1)', 'requisition', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(276, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(277, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(278, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(279, 'INFO', '编辑角色: 项目经理', 'auth', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(280, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-06-30 00:00:00'),
(281, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(282, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(283, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(284, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(285, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-06-30 00:00:00'),
(286, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-07-01 00:00:00'),
(287, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-01 00:00:00'),
(288, 'INFO', '创建入库单: RK20250702081416', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(289, 'INFO', '创建入库单: RK20250702081717', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(290, 'INFO', '创建入库单: RK20250702081808', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(291, 'INFO', '创建入库单: RK20250702081934', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(292, 'INFO', '创建入库单: RK20250702082034', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(293, 'INFO', '创建入库单: RK20250702082120', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(294, 'INFO', '创建入库单: RK20250702082213', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(295, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(296, 'INFO', '新增物资: 镀锌钢管 (编码:None)', 'material', 1, '172.18.0.1', '2025-07-02 00:00:00');
,
INSERT INTO "system_logs" ("id", "level", "message", "module", "user_id", "ip_address", "created_at") VALUES
(297, 'INFO', '编辑物资: 镀锌钢管 (编码:None)', 'material', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(298, 'INFO', '新增物资: 电话电缆 (编码:None)', 'material', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(299, 'INFO', '创建入库单: RK20250702144627', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(300, 'INFO', '审批入库单: RK20250702081416，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(301, 'INFO', '审批入库单: RK20250702081717，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(302, 'INFO', '审批入库单: RK20250702081808，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(303, 'INFO', '审批入库单: RK20250702081934，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(304, 'INFO', '审批入库单: RK20250702082034，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(305, 'INFO', '审批入库单: RK20250702082120，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(306, 'INFO', '审批入库单: RK20250702082213，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(307, 'INFO', '审批入库单: RK20250702144627，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(308, 'INFO', '创建入库单: RK20250702145258', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(309, 'INFO', '创建入库单: RK20250702145437', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(310, 'INFO', '创建入库单: RK20250702145614', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(311, 'INFO', '审批入库单: RK20250702145258，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(312, 'INFO', '审批入库单: RK20250702145437，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(313, 'INFO', '审批入库单: RK20250702145614，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(314, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-07-02 00:00:00'),
(315, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(316, 'INFO', '新增申请单: CK20250702001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(317, 'INFO', '新增申请单: CK20250702002 (项目ID:2)', 'requisition', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(318, 'INFO', '新增申请单: CK20250702003 (项目ID:9)', 'requisition', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(319, 'WARNING', '删除申请单: CK20250702001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(320, 'INFO', '新增申请单: CK20250702003 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(321, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-07-02 00:00:00'),
(322, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(323, 'INFO', '审核申请单: CK20250702003 (项目ID:1)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(324, 'INFO', '审核申请单: CK20250702003 (项目ID:9)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(325, 'INFO', '审核申请单: CK20250702002 (项目ID:2)，结果: 通过', 'requisition', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(326, 'INFO', '发放申请单: CK20250702002 (项目ID:2)', 'requisition', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(327, 'INFO', '发放申请单: CK20250702003 (项目ID:9)', 'requisition', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(328, 'INFO', '发放申请单: CK20250702003 (项目ID:1)', 'requisition', 3, '172.18.0.1', '2025-07-02 00:00:00'),
(329, 'INFO', '编辑项目: 腈纶厂长丝项目仪表附属设备区 (编号:PJ-20250622001)', 'project', 3, '172.18.0.1', '2025-07-03 00:00:00'),
(330, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-07-03 00:00:00'),
(331, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-03 00:00:00'),
(332, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-05 00:00:00'),
(333, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-06 00:00:00'),
(334, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-07 00:00:00'),
(335, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-09 00:00:00'),
(336, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-10 00:00:00'),
(337, 'INFO', '创建入库单: RK20250711040424', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(338, 'INFO', '创建入库单: RK20250711040505', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(339, 'INFO', '创建入库单: RK20250711040525', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(340, 'INFO', '审批入库单: RK20250711040525，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(341, 'INFO', '审批入库单: RK20250711040505，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(342, 'INFO', '审批入库单: RK20250711040424，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(343, 'INFO', '创建入库单: RK20250711041036', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(344, 'INFO', '审批入库单: RK20250711041036，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(345, 'INFO', '创建入库单: RK20250711041508', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(346, 'INFO', '创建入库单: RK20250711041828', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(347, 'INFO', '审批入库单: RK20250711041508，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(348, 'INFO', '审批入库单: RK20250711041828，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(349, 'INFO', '创建入库单: RK20250711042020', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(350, 'INFO', '创建入库单: RK20250711042055', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(351, 'INFO', '创建入库单: RK20250711042119', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(352, 'INFO', '创建入库单: RK20250711042142', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(353, 'INFO', '审批入库单: RK20250711042020，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(354, 'INFO', '审批入库单: RK20250711042055，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(355, 'INFO', '审批入库单: RK20250711042119，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(356, 'INFO', '审批入库单: RK20250711042142，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(357, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(358, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(359, 'INFO', '新增申请单: CK20250711001 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(360, 'INFO', '新增申请单: CK20250711002 (项目ID:2)', 'requisition', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(361, 'INFO', '新增申请单: CK20250711003 (项目ID:4)', 'requisition', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(362, 'INFO', '新增申请单: CK20250711004 (项目ID:5)', 'requisition', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(363, 'INFO', '新增申请单: CK20250711005 (项目ID:6)', 'requisition', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(364, 'INFO', '新增申请单: CK20250711006 (项目ID:9)', 'requisition', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(365, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-07-11 00:00:00'),
(366, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(367, 'INFO', '审核申请单: CK20250711001 (项目ID:1)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(368, 'INFO', '发放申请单: CK20250711001 (项目ID:1)', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(369, 'INFO', '审核申请单: CK20250711002 (项目ID:2)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(370, 'INFO', '审核申请单: CK20250711003 (项目ID:4)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(371, 'INFO', '审核申请单: CK20250711004 (项目ID:5)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(372, 'INFO', '审核申请单: CK20250711005 (项目ID:6)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(373, 'INFO', '审核申请单: CK20250711006 (项目ID:9)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(374, 'INFO', '发放申请单: CK20250711002 (项目ID:2)', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(375, 'INFO', '发放申请单: CK20250711003 (项目ID:4)', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(376, 'INFO', '发放申请单: CK20250711004 (项目ID:5)', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(377, 'INFO', '发放申请单: CK20250711005 (项目ID:6)', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(378, 'INFO', '发放申请单: CK20250711006 (项目ID:9)', 'requisition', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(379, 'WARNING', '删除项目: 腈纶厂长丝项目电信 (编号:PJ-20250622003)', 'project', 1, '172.18.0.1', '2025-07-11 00:00:00'),
(380, 'INFO', '编辑项目: 腈纶厂长丝项目火灾 (编号:PJ-20250622004)', 'project', 1, '172.18.0.1', '2025-07-14 00:00:00'),
(381, 'INFO', '编辑项目: 腈纶厂长丝项目视频监控 (编号:PJ-20250622005)', 'project', 1, '172.18.0.1', '2025-07-14 00:00:00'),
(382, 'INFO', '编辑项目: 腈纶厂长丝项目广播系统 (编号:PJ-20250622006)', 'project', 1, '172.18.0.1', '2025-07-14 00:00:00'),
(383, 'INFO', '编辑项目: 腈纶厂长丝项目门禁系统 (编号:PJ-20250622007)', 'project', 1, '172.18.0.1', '2025-07-14 00:00:00'),
(384, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 1, '172.18.0.1', '2025-07-14 00:00:00'),
(385, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 1, '172.18.0.1', '2025-07-14 00:00:00'),
(386, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-17 00:00:00'),
(387, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-21 00:00:00'),
(388, 'INFO', '更新系统设置: [''system_name'', ''session_timeout'']', 'system', 1, '172.18.0.1', '2025-07-21 00:00:00'),
(389, 'INFO', '执行数据库备份: backup_20250721_003210.sql', 'backup', 1, '172.18.0.1', '2025-07-21 00:00:00'),
(390, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-21 00:00:00'),
(391, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(392, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(393, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(394, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(395, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(396, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00');
,
INSERT INTO "system_logs" ("id", "level", "message", "module", "user_id", "ip_address", "created_at") VALUES
(397, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(398, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(399, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(400, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(401, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(402, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(403, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(404, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(405, 'WARNING', '删除物资: 膨胀螺栓(配套螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(406, 'WARNING', '删除物资: 膨胀螺栓(配套螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(407, 'WARNING', '删除物资: 膨胀螺栓(配套螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(408, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(409, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(410, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(411, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(412, 'WARNING', '删除物资: U型管卡(带锁紧螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(413, 'WARNING', '删除物资: 膨胀螺栓(配套螺母) (编码:None)', 'material', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(414, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(415, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(416, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(417, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(418, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(419, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(420, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(421, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(422, 'INFO', '编辑角色: 分包材料员', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(423, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(424, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(425, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-07-22 00:00:00'),
(426, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(427, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-22 00:00:00'),
(428, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-07-23 00:00:00'),
(429, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-23 00:00:00'),
(430, 'INFO', '新增项目: 车间库存 (编号:PJ-20250723001)', 'project', 1, '172.18.0.1', '2025-07-23 00:00:00'),
(431, 'INFO', '编辑项目: 车间库存 (编号:PJ-20250723001)', 'project', 1, '172.18.0.1', '2025-07-25 00:00:00'),
(432, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-25 00:00:00'),
(433, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-30 00:00:00'),
(434, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-30 00:00:00'),
(435, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-07-30 00:00:00'),
(436, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(437, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(438, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(439, 'INFO', '新增申请单: CK20250804001 (项目ID:9)', 'requisition', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(440, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(441, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(442, 'INFO', '创建入库单: RK20250804011947', 'inbound', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(443, 'INFO', '审批入库单: RK20250804011947，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(444, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(445, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(446, 'INFO', '新增申请单: CK20250804002 (项目ID:1)', 'requisition', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(447, 'INFO', '新增申请单: CK20250804003 (项目ID:2)', 'requisition', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(448, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-08-04 00:00:00'),
(449, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(450, 'INFO', '审核申请单: CK20250804003 (项目ID:2)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(451, 'INFO', '审核申请单: CK20250804002 (项目ID:1)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(452, 'INFO', '审核申请单: CK20250804001 (项目ID:9)，结果: 通过', 'requisition', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(453, 'INFO', '发放申请单: CK20250804003 (项目ID:2)', 'requisition', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(454, 'INFO', '发放申请单: CK20250804002 (项目ID:1)', 'requisition', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(455, 'INFO', '发放申请单: CK20250804001 (项目ID:9)', 'requisition', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(456, 'INFO', '库存调整: 镀锌角钢 (类型:decrease, 新库存:0.0)', 'stock', 1, '172.18.0.1', '2025-08-04 00:00:00'),
(457, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-08-05 00:00:00'),
(458, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-08-07 00:00:00'),
(459, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-08-22 00:00:00'),
(460, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(461, 'INFO', '创建入库单: RK20250905012310', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(462, 'INFO', '创建入库单: RK20250905012404', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(463, 'INFO', '创建入库单: RK20250905012435', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(464, 'INFO', '创建入库单: RK20250905012502', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(465, 'INFO', '创建入库单: RK20250905012534', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(466, 'INFO', '创建入库单: RK20250905012550', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(467, 'INFO', '创建入库单: RK20250905012608', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(468, 'INFO', '创建入库单: RK20250905012630', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(469, 'INFO', '审批入库单: RK20250905012630，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(470, 'INFO', '审批入库单: RK20250905012608，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(471, 'INFO', '审批入库单: RK20250905012550，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(472, 'INFO', '审批入库单: RK20250905012534，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(473, 'INFO', '审批入库单: RK20250905012502，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(474, 'INFO', '审批入库单: RK20250905012435，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(475, 'INFO', '审批入库单: RK20250905012404，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(476, 'INFO', '审批入库单: RK20250905012310，结果: 批准', 'inbound', 1, '172.18.0.1', '2025-09-05 00:00:00'),
(477, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-08 00:00:00'),
(478, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-08 00:00:00'),
(479, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-09 00:00:00'),
(480, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-10 00:00:00'),
(481, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(482, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(483, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(484, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(485, 'INFO', '编辑项目: 车间库存 (编号:PJ-20250723001)', 'project', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(486, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(487, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(488, 'INFO', '编辑项目: 车间库存 (编号:PJ-20250723001)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(489, 'WARNING', '删除项目: 车间库存 (编号:PJ-20250723001)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(490, 'INFO', '新增项目: test (编号:PJ-20250911001)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(491, 'INFO', '新增物资: 213123 (编码:None)', 'material', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(492, 'WARNING', '删除物资: 213123 (编码:None)', 'material', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(493, 'WARNING', '删除项目: test (编号:PJ-20250911001)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(494, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(495, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(496, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00');
INSERT INTO "system_logs" ("id", "level", "message", "module", "user_id", "ip_address", "created_at") VALUES
(497, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(498, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-09-11 00:00:00'),
(499, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(500, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(501, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 1, '172.18.0.1', '2025-09-11 00:00:00'),
(502, 'INFO', '用户登出', 'auth', 1, '172.18.0.1', '2025-09-12 00:00:00'),
(503, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-09-12 00:00:00'),
(504, 'INFO', '用户登出', 'auth', 3, '172.18.0.1', '2025-09-12 00:00:00'),
(505, 'INFO', '用户登录: gjw', 'auth', 5, '172.18.0.1', '2025-09-12 00:00:00'),
(506, 'INFO', '用户登出', 'auth', 5, '172.18.0.1', '2025-09-12 00:00:00'),
(507, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-09-12 00:00:00'),
(508, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-09-22 00:00:00'),
(509, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-10-09 00:00:00'),
(510, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-10-09 00:00:00'),
(511, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-10-28 00:00:00'),
(512, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(513, 'INFO', '编辑项目: 腈纶厂长丝项目仪表附属设备区 (编号:PJ-20250622001)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(514, 'INFO', '编辑项目: 腈纶厂长丝项目仪表专有设备区 (编号:PJ-20250622002)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(515, 'INFO', '编辑项目: 腈纶厂长丝项目火灾 (编号:PJ-20250622004)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(516, 'INFO', '编辑项目: 腈纶厂长丝项目视频监控 (编号:PJ-20250622005)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(517, 'INFO', '编辑项目: 腈纶厂长丝项目广播系统 (编号:PJ-20250622006)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(518, 'INFO', '编辑项目: 腈纶厂长丝项目门禁系统 (编号:PJ-20250622007)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(519, 'INFO', '编辑项目: 腈纶厂长丝项目电话系统 (编号:PJ-20250622008)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(520, 'INFO', '编辑项目: 腈纶厂长丝项目空调机组 (编号:PJ-20250622009)', 'project', 3, '172.18.0.1', '2025-11-06 00:00:00'),
(521, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-11-17 00:00:00'),
(522, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-12-16 00:00:00'),
(523, 'INFO', '执行数据库备份: backup_20251216_070717.sql', 'backup', 1, '172.18.0.1', '2025-12-16 00:00:00'),
(524, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-12-16 00:00:00'),
(525, 'INFO', '用户登录: julei', 'auth', 3, '172.18.0.1', '2025-12-17 00:00:00'),
(526, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-12-26 00:00:00'),
(527, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2025-12-26 00:00:00'),
(528, 'INFO', '用户登录: admin', 'auth', 1, '172.18.0.1', '2026-01-21 00:00:00'),
(529, 'INFO', '新增项目: 线性装置聚烯烃弹性体生产优化改造 (编号:PJ-20260121001)', 'project', 1, '172.18.0.1', '2026-01-21 00:00:00');

-- Table: task_dependencies

-- Schema and data for table: task_dependencies
INSERT INTO "task_dependencies" ("id", "task_id", "depends_on", "type", "lag", "created_at") VALUES
(2, 6, 5, 'FS', 0, '2026-01-31 08:28:50'),
(3, 4, 3, 'FS', 0, '2026-01-31 08:32:28'),
(4, 5, 4, 'FS', 0, '2026-01-31 08:32:34'),
(5, 7, 6, 'FS', 0, '2026-01-31 08:38:37');

-- Table: task_resources

-- Schema and data for table: task_resources
INSERT INTO "task_resources" ("id", "task_id", "resource_id", "quantity", "created_at", "updated_at") VALUES
(1, 3, 3, '1', '2026-01-31 13:05:50', '2026-01-31 13:08:08'),
(2, 3, 4, '1', '2026-01-31 13:08:08', '2026-01-31 13:08:08'),
(3, 4, 4, '1', '2026-01-31 13:35:17', '2026-01-31 13:35:40'),
(4, 4, 3, '1', '2026-01-31 13:35:40', '2026-01-31 13:35:40');

-- Table: tasks

-- Schema and data for table: tasks
INSERT INTO "tasks" ("id", "project_id", "schedule_id", "parent_id", "name", "duration", "start_date", "end_date", "progress", "is_milestone", "sort_order", "position_x", "position_y", "created_at", "updated_at", "priority", "status", "responsible", "description") VALUES
(1, 9, NULL, NULL, '测试任务1', '0.00', '2026-01-30 00:00:00', '2026-01-30 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 01:54:00', '2026-01-30 02:26:43', 'medium', 'not_started', '', ''),
(7, 10, NULL, 5, '施工准备', '80.00', '2026-01-22 00:00:00', '2026-04-12 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:27:08', '2026-02-01 01:00:45', 'medium', 'not_started', '', ''),
(9, 10, NULL, 3, '质量检查', '7.00', '2026-04-05 00:00:00', '2026-04-12 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:27:08', '2026-02-01 01:00:45', 'medium', 'not_started', '', ''),
(8, 10, NULL, 5, '施工执行', '7.00', '2026-04-05 00:00:00', '2026-04-12 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:27:08', '2026-02-01 01:00:45', 'medium', 'not_started', '', ''),
(11, 10, NULL, 10, 'ttt', '3.00', '2026-01-22 00:00:00', '2026-01-25 00:00:00', '0.00', true, 0, NULL, NULL, '2026-01-30 05:10:09', '2026-01-31 15:18:03', 'medium', 'not_started', '', ''),
(3, 10, NULL, 4, '测试1', '80.00', '2026-01-22 00:00:00', '2026-04-12 00:00:00', '7.00', false, 0, NULL, NULL, '2026-01-30 02:44:13', '2026-02-02 04:25:54', 'medium', 'not_started', '', ''),
(5, 10, NULL, 4, '项目启动', '80.00', '2026-01-22 00:00:00', '2026-04-12 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:27:08', '2026-02-02 04:26:03', 'medium', 'not_started', '', ''),
(10, 10, NULL, 3, '竣工验收', '3.00', '2026-01-22 00:00:00', '2026-01-25 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:27:08', '2026-01-31 15:25:10', 'medium', 'not_started', '', ''),
(4, 10, NULL, 3, '项目启动', '80.00', '2026-01-22 00:00:00', '2026-04-12 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:26:49', '2026-02-02 04:26:03', 'medium', 'not_started', '', ''),
(6, 10, NULL, 5, '方案设计', '80.00', '2026-01-22 00:00:00', '2026-04-12 00:00:00', '0.00', false, 0, NULL, NULL, '2026-01-30 04:27:08', '2026-02-01 01:00:45', 'medium', 'not_started', '', '');

-- Table: user_projects

-- Schema and data for table: user_projects
INSERT INTO "user_projects" ("user_id", "project_id") VALUES
(3, 10),
(4, 10),
(1, 10),
(3, 1),
(2, 1),
(1, 1),
(1, 2),
(1, 4),
(3, 4),
(3, 2),
(2, 2),
(2, 5),
(3, 5),
(3, 6),
(2, 6),
(3, 7),
(2, 7),
(2, 8),
(5, 1),
(5, 2),
(4, 4),
(5, 4),
(4, 5),
(5, 5),
(4, 6),
(5, 6),
(5, 8),
(4, 8),
(4, 1),
(3, 8),
(5, 11),
(4, 11),
(3, 11),
(4, 9),
(3, 9);

-- Table: user_roles

-- Schema and data for table: user_roles
INSERT INTO "user_roles" ("user_id", "role_id", "id") VALUES
(1, 1, 1),
(2, 2, 2),
(3, 5, 3),
(4, 5, 4),
(5, 4, 5),
(3, 3, 6),
(1, 5, 7);

-- Table: users

-- Schema and data for table: users
INSERT INTO "users" ("id", "username", "password", "email", "full_name", "role", "group", "is_active", "last_login", "created_at") VALUES
(3, 'julei', 'scrypt:32768:8:1$6N27g9ASq7GPgpFv$0ba3a928b40c699c1818ea6f98874c5994222a6e84b42226f2db4c4fcf96f3bc919876f81c16bcd80d7b7ebdb6a958ed6144c136581f8106d3383774c7df50e0', 'cnjulei@gmail.com', '鞠磊', '项目经理', '', true, '2026-02-03 00:00:00', '2025-06-22 00:00:00'),
(5, 'gjw', 'scrypt:32768:8:1$xR9riv8ulBQehmfj$aeb45718560cc94eac78582426b97a888ce3b6c50f104de096a2569888cb1a1e0d5b9b89124b168a98fadac288bf289c1d92f0edc9f92ef7ae76dc74f06f53a4', 'gjw@126.com', '葛继伟', '分包材料员', '', true, '2026-01-26 00:00:00', '2025-06-28 00:00:00'),
(4, 'libo', 'scrypt:32768:8:1$AsKliUJQzeQ71ax4$5a1b727c8335316e9c7e8251320389e5c2e364b8544930bd7123c598c51b9de37e4ab1045c3cf9d631d31d2af99d9a2b3c64ec48cb3c69aee5667b2414774c15', 'libo@166.com', '李波', '项目经理', '', true, '2026-01-27 00:00:00', '2025-06-28 00:00:00'),
(2, 'wqs', 'scrypt:32768:8:1$vRD3m50BevQLxI32$bd76ee613eba79b52feba85d3c4439abed51d0f06c44e4a8490195bc3b6b2114316e09efb98d7d4960f1dbc8ac2a755878306897c368d92fb78c1b05a7449fd2', 'test@gmail.com', '王庆森', '保管员', '', true, '2026-02-02 00:00:00', '2025-06-22 00:00:00'),
(1, 'admin', '$2a$10$7caGzyrKOa5SVKgVIUeHdOfOrUEwWUhrL/RhxtfNquS6np912C/s.', 'admin@example.com', '系统管理员', 'admin', '', true, '2026-02-04 00:00:00', '2025-06-21 00:00:00');

-- Table: workflow_approvals

-- Schema and data for table: workflow_approvals
INSERT INTO "workflow_approvals" ("id", "instance_id", "node_id", "node_key", "approver_id", "approver_name", "action", "remark", "attachments", "approved_at", "created_at") VALUES
(1, 1, 94, 'node_1', 1, 'admin', 'approve', '', '', '2026-02-03 08:51:43', '2026-02-03 08:51:43'),
(2, 2, 100, 'node_1', 1, 'admin', 'approve', '', '', '2026-02-03 09:01:30', '2026-02-03 09:01:30'),
(3, 3, 100, 'node_1', 1, 'admin', 'approve', '测试审核通过', '', '2026-02-04 05:42:50', '2026-02-04 05:42:50'),
(4, 4, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 13:37:27', '2026-02-04 13:37:27'),
(5, 5, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 13:38:22', '2026-02-04 13:38:22'),
(6, 6, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 13:40:03', '2026-02-04 13:40:03'),
(7, 7, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 13:45:08', '2026-02-04 13:45:08'),
(8, 8, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 13:50:16', '2026-02-04 13:50:16'),
(9, 9, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 13:54:07', '2026-02-04 13:54:07'),
(10, 10, 97, 'node_1', 1, '未知用户', 'approve', '', '', '2026-02-04 14:05:01', '2026-02-04 14:05:01');

-- Table: workflow_definitions

-- Schema and data for table: workflow_definitions
INSERT INTO "workflow_definitions" ("id", "name", "description", "module", "version", "is_active", "created_at", "updated_at") VALUES
(7, '计划审批工作流', '计划审批工作流', 'material_plan', 1, true, '2026-02-02 04:09:15', '2026-02-02 04:09:15'),
(8, '入库管理工作流', '入库管理工作流', 'inbound', 1, true, '2026-02-02 04:09:41', '2026-02-02 04:09:41'),
(9, '出库管理工作流', '出库管理工作流', 'requisition', 1, true, '2026-02-02 04:10:18', '2026-02-02 04:10:18');

-- Table: workflow_edges

-- Schema and data for table: workflow_edges
INSERT INTO "workflow_edges" ("id", "workflow_id", "from_node", "to_node", "condition_expression", "created_at") VALUES
(68, 9, 'start', 'node_1', '', '2026-02-02 04:10:18'),
(69, 9, 'node_1', 'end', '', '2026-02-02 04:10:18'),
(66, 8, 'start', 'node_1', '', '2026-02-02 04:09:41'),
(67, 8, 'node_1', 'end', '', '2026-02-02 04:09:41'),
(64, 7, 'start', 'node_1', '', '2026-02-02 04:09:15'),
(65, 7, 'node_1', 'end', '', '2026-02-02 04:09:15');

-- Table: workflow_instances

-- Schema and data for table: workflow_instances
INSERT INTO "workflow_instances" ("id", "workflow_id", "business_type", "business_id", "business_no", "current_node", "status", "initiator_id", "initiator_name", "started_at", "finished_at", "created_at", "updated_at") VALUES
(4, 8, 'inbound_order', 11, 'RK20260204133723', 'end', 'approved', 1, 'admin', '2026-02-04 13:37:23', '2026-02-04 13:37:27', '2026-02-04 13:37:23', '2026-02-04 13:37:27'),
(5, 8, 'inbound_order', 12, 'RK20260204133807', 'end', 'approved', 1, 'admin', '2026-02-04 13:38:07', '2026-02-04 13:38:22', '2026-02-04 13:38:07', '2026-02-04 13:38:22'),
(6, 8, 'inbound_order', 13, 'RK20260204133954', 'end', 'approved', 1, 'admin', '2026-02-04 13:39:55', '2026-02-04 13:40:03', '2026-02-04 13:39:55', '2026-02-04 13:40:03'),
(7, 8, 'inbound_order', 14, 'RK20260204134503', 'end', 'approved', 1, 'admin', '2026-02-04 13:45:03', '2026-02-04 13:45:08', '2026-02-04 13:45:03', '2026-02-04 13:45:08'),
(8, 8, 'inbound_order', 15, 'RK20260204135013', 'end', 'approved', 1, 'admin', '2026-02-04 13:50:13', '2026-02-04 13:50:16', '2026-02-04 13:50:13', '2026-02-04 13:50:16'),
(9, 8, 'inbound_order', 16, 'RK20260204135403', 'end', 'approved', 1, 'admin', '2026-02-04 13:54:03', '2026-02-04 13:54:07', '2026-02-04 13:54:03', '2026-02-04 13:54:07'),
(10, 8, 'inbound_order', 17, 'RK20260204140457', 'end', 'approved', 1, 'admin', '2026-02-04 14:04:57', '2026-02-04 14:05:01', '2026-02-04 14:04:57', '2026-02-04 14:05:01'),
(1, 7, 'material_plan', 1, 'MP2602030001', 'end', 'approved', 1, 'admin', '2026-02-03 08:51:37', '2026-02-03 08:51:43', '2026-02-03 08:51:37', '2026-02-03 08:51:43'),
(2, 9, 'requisition', 1, 'CK20260203001', 'end', 'approved', 1, 'admin', '2026-02-03 09:01:27', '2026-02-03 09:01:30', '2026-02-03 09:01:27', '2026-02-03 09:01:30'),
(3, 9, 'requisition', 2, 'CK20260204001', 'end', 'approved', 1, 'admin', '2026-02-04 05:42:43', '2026-02-04 05:42:50', '2026-02-04 05:42:43', '2026-02-04 05:42:50');

-- Table: workflow_logs

-- Schema and data for table: workflow_logs
INSERT INTO "workflow_logs" ("id", "instance_id", "node_key", "action", "actor_id", "actor_name", "action_data", "created_at") VALUES
(102, 1, 'start', 'start', 1, 'admin', '', '2026-02-03 08:51:37'),
(103, 1, 'node_1', 'approve', 1, 'admin', '', '2026-02-03 08:51:43'),
(104, 1, 'end', 'finish', 1, 'admin', '', '2026-02-03 08:51:43'),
(105, 2, 'start', 'start', 1, 'admin', '', '2026-02-03 09:01:27'),
(106, 2, 'node_1', 'approve', 1, 'admin', '', '2026-02-03 09:01:30'),
(107, 2, 'end', 'finish', 1, 'admin', '', '2026-02-03 09:01:30'),
(108, 3, 'start', 'start', 1, 'admin', '', '2026-02-04 05:42:43'),
(109, 3, 'node_1', 'approve', 1, 'admin', '', '2026-02-04 05:42:50'),
(110, 3, 'end', 'finish', 1, 'admin', '', '2026-02-04 05:42:50'),
(111, 4, 'start', 'start', 1, 'admin', '', '2026-02-04 13:37:23'),
(112, 4, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 13:37:27'),
(113, 4, 'end', 'finish', 1, '未知用户', '', '2026-02-04 13:37:27'),
(114, 5, 'start', 'start', 1, 'admin', '', '2026-02-04 13:38:07'),
(115, 5, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 13:38:22'),
(116, 5, 'end', 'finish', 1, '未知用户', '', '2026-02-04 13:38:22'),
(117, 6, 'start', 'start', 1, 'admin', '', '2026-02-04 13:39:55'),
(118, 6, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 13:40:03'),
(119, 6, 'end', 'finish', 1, '未知用户', '', '2026-02-04 13:40:03'),
(120, 7, 'start', 'start', 1, 'admin', '', '2026-02-04 13:45:03'),
(121, 7, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 13:45:08'),
(122, 7, 'end', 'finish', 1, '未知用户', '', '2026-02-04 13:45:08'),
(123, 8, 'start', 'start', 1, 'admin', '', '2026-02-04 13:50:13'),
(124, 8, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 13:50:16'),
(125, 8, 'end', 'finish', 1, '未知用户', '', '2026-02-04 13:50:16'),
(126, 9, 'start', 'start', 1, 'admin', '', '2026-02-04 13:54:03'),
(127, 9, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 13:54:07'),
(128, 9, 'end', 'finish', 1, '未知用户', '', '2026-02-04 13:54:07'),
(129, 10, 'start', 'start', 1, 'admin', '', '2026-02-04 14:04:57'),
(130, 10, 'node_1', 'approve', 1, '未知用户', '', '2026-02-04 14:05:01'),
(131, 10, 'end', 'finish', 1, '未知用户', '', '2026-02-04 14:05:01');

-- Table: workflow_node_approvers

-- Schema and data for table: workflow_node_approvers
INSERT INTO "workflow_node_approvers" ("id", "node_id", "approver_type", "approver_id", "approver_name", "sequence", "created_at") VALUES
(30, 94, 'role', 5, '项目经理', 0, '2026-02-02 04:09:15'),
(31, 97, 'role', 5, '项目经理', 0, '2026-02-02 04:09:41'),
(32, 100, 'role', 5, '项目经理', 0, '2026-02-02 04:10:18');

-- Table: workflow_nodes

-- Schema and data for table: workflow_nodes
INSERT INTO "workflow_nodes" ("id", "workflow_id", "node_key", "node_type", "node_name", "description", "approval_type", "timeout_hours", "auto_approve", "is_required", "created_at", "x", "y") VALUES
(92, 7, 'start', 'start', '开始', '流程开始', 'sequential', 0, false, true, '2026-02-02 04:09:15', 100, 200),
(93, 7, 'end', 'end', '结束', '流程结束', 'sequential', 0, false, true, '2026-02-02 04:09:15', 500, 200),
(94, 7, 'node_1', 'approval', '新节点', '', 'sequential', 24, false, true, '2026-02-02 04:09:15', 280, 189),
(95, 8, 'start', 'start', '开始', '流程开始', 'sequential', 0, false, true, '2026-02-02 04:09:41', 100, 200),
(96, 8, 'end', 'end', '结束', '流程结束', 'sequential', 0, false, true, '2026-02-02 04:09:41', 500, 200),
(97, 8, 'node_1', 'approval', '新节点', '', 'sequential', 24, false, true, '2026-02-02 04:09:41', 265, 184),
(98, 9, 'start', 'start', '开始', '流程开始', 'sequential', 0, false, true, '2026-02-02 04:10:18', 100, 200),
(99, 9, 'end', 'end', '结束', '流程结束', 'sequential', 0, false, true, '2026-02-02 04:10:18', 500, 200),
(100, 9, 'node_1', 'approval', '新节点', '', 'sequential', 24, false, true, '2026-02-02 04:10:18', 263, 195);

-- Table: workflow_pending_tasks

-- Schema and data for table: workflow_pending_tasks
INSERT INTO "workflow_pending_tasks" ("id", "instance_id", "node_id", "node_key", "node_name", "business_type", "business_id", "business_no", "approver_id", "approver_name", "status", "is_parallel", "arrived_at", "processed_at", "created_at", "updated_at") VALUES
(121, 9, 97, 'node_1', '新节点', 'inbound_order', 16, 'RK20260204135403', 3, '鞠磊', 'pending', false, '2026-02-04 13:54:03', NULL, '2026-02-04 13:54:03', '2026-02-04 13:54:03'),
(122, 9, 97, 'node_1', '新节点', 'inbound_order', 16, 'RK20260204135403', 4, '李波', 'pending', false, '2026-02-04 13:54:03', NULL, '2026-02-04 13:54:03', '2026-02-04 13:54:03'),
(120, 9, 97, 'node_1', '新节点', 'inbound_order', 16, 'RK20260204135403', 1, '系统管理员', 'approved', false, '2026-02-04 13:54:03', '2026-02-04 13:54:07', '2026-02-04 13:54:03', '2026-02-04 13:54:07'),
(124, 10, 97, 'node_1', '新节点', 'inbound_order', 17, 'RK20260204140457', 3, '鞠磊', 'pending', false, '2026-02-04 14:04:58', NULL, '2026-02-04 14:04:58', '2026-02-04 14:04:58'),
(125, 10, 97, 'node_1', '新节点', 'inbound_order', 17, 'RK20260204140457', 4, '李波', 'pending', false, '2026-02-04 14:04:58', NULL, '2026-02-04 14:04:58', '2026-02-04 14:04:58'),
(123, 10, 97, 'node_1', '新节点', 'inbound_order', 17, 'RK20260204140457', 1, '系统管理员', 'approved', false, '2026-02-04 14:04:58', '2026-02-04 14:05:01', '2026-02-04 14:04:58', '2026-02-04 14:05:01'),
(97, 1, 94, 'node_1', '新节点', 'material_plan', 1, 'MP2602030001', 3, '鞠磊', 'pending', false, '2026-02-03 08:51:37', NULL, '2026-02-03 08:51:37', '2026-02-03 08:51:37'),
(98, 1, 94, 'node_1', '新节点', 'material_plan', 1, 'MP2602030001', 4, '李波', 'pending', false, '2026-02-03 08:51:37', NULL, '2026-02-03 08:51:37', '2026-02-03 08:51:37'),
(96, 1, 94, 'node_1', '新节点', 'material_plan', 1, 'MP2602030001', 1, '系统管理员', 'approved', false, '2026-02-03 08:51:37', '2026-02-03 08:51:43', '2026-02-03 08:51:37', '2026-02-03 08:51:43'),
(100, 2, 100, 'node_1', '新节点', 'requisition', 1, 'CK20260203001', 3, '鞠磊', 'pending', false, '2026-02-03 09:01:27', NULL, '2026-02-03 09:01:27', '2026-02-03 09:01:27'),
(101, 2, 100, 'node_1', '新节点', 'requisition', 1, 'CK20260203001', 4, '李波', 'pending', false, '2026-02-03 09:01:27', NULL, '2026-02-03 09:01:27', '2026-02-03 09:01:27'),
(99, 2, 100, 'node_1', '新节点', 'requisition', 1, 'CK20260203001', 1, '系统管理员', 'approved', false, '2026-02-03 09:01:27', '2026-02-03 09:01:30', '2026-02-03 09:01:27', '2026-02-03 09:01:30'),
(103, 3, 100, 'node_1', '新节点', 'requisition', 2, 'CK20260204001', 3, '鞠磊', 'pending', false, '2026-02-04 05:42:43', NULL, '2026-02-04 05:42:43', '2026-02-04 05:42:43'),
(104, 3, 100, 'node_1', '新节点', 'requisition', 2, 'CK20260204001', 4, '李波', 'pending', false, '2026-02-04 05:42:43', NULL, '2026-02-04 05:42:43', '2026-02-04 05:42:43'),
(102, 3, 100, 'node_1', '新节点', 'requisition', 2, 'CK20260204001', 1, '系统管理员', 'approved', false, '2026-02-04 05:42:43', '2026-02-04 05:42:50', '2026-02-04 05:42:43', '2026-02-04 05:42:50'),
(106, 4, 97, 'node_1', '新节点', 'inbound_order', 11, 'RK20260204133723', 3, '鞠磊', 'pending', false, '2026-02-04 13:37:23', NULL, '2026-02-04 13:37:23', '2026-02-04 13:37:23'),
(107, 4, 97, 'node_1', '新节点', 'inbound_order', 11, 'RK20260204133723', 4, '李波', 'pending', false, '2026-02-04 13:37:23', NULL, '2026-02-04 13:37:23', '2026-02-04 13:37:23'),
(105, 4, 97, 'node_1', '新节点', 'inbound_order', 11, 'RK20260204133723', 1, '系统管理员', 'approved', false, '2026-02-04 13:37:23', '2026-02-04 13:37:27', '2026-02-04 13:37:23', '2026-02-04 13:37:27'),
(109, 5, 97, 'node_1', '新节点', 'inbound_order', 12, 'RK20260204133807', 3, '鞠磊', 'pending', false, '2026-02-04 13:38:07', NULL, '2026-02-04 13:38:07', '2026-02-04 13:38:07'),
(110, 5, 97, 'node_1', '新节点', 'inbound_order', 12, 'RK20260204133807', 4, '李波', 'pending', false, '2026-02-04 13:38:07', NULL, '2026-02-04 13:38:07', '2026-02-04 13:38:07'),
(108, 5, 97, 'node_1', '新节点', 'inbound_order', 12, 'RK20260204133807', 1, '系统管理员', 'approved', false, '2026-02-04 13:38:07', '2026-02-04 13:38:22', '2026-02-04 13:38:07', '2026-02-04 13:38:22'),
(112, 6, 97, 'node_1', '新节点', 'inbound_order', 13, 'RK20260204133954', 3, '鞠磊', 'pending', false, '2026-02-04 13:39:55', NULL, '2026-02-04 13:39:55', '2026-02-04 13:39:55'),
(113, 6, 97, 'node_1', '新节点', 'inbound_order', 13, 'RK20260204133954', 4, '李波', 'pending', false, '2026-02-04 13:39:55', NULL, '2026-02-04 13:39:55', '2026-02-04 13:39:55'),
(111, 6, 97, 'node_1', '新节点', 'inbound_order', 13, 'RK20260204133954', 1, '系统管理员', 'approved', false, '2026-02-04 13:39:55', '2026-02-04 13:40:03', '2026-02-04 13:39:55', '2026-02-04 13:40:03'),
(115, 7, 97, 'node_1', '新节点', 'inbound_order', 14, 'RK20260204134503', 3, '鞠磊', 'pending', false, '2026-02-04 13:45:03', NULL, '2026-02-04 13:45:03', '2026-02-04 13:45:03'),
(116, 7, 97, 'node_1', '新节点', 'inbound_order', 14, 'RK20260204134503', 4, '李波', 'pending', false, '2026-02-04 13:45:03', NULL, '2026-02-04 13:45:03', '2026-02-04 13:45:03'),
(114, 7, 97, 'node_1', '新节点', 'inbound_order', 14, 'RK20260204134503', 1, '系统管理员', 'approved', false, '2026-02-04 13:45:03', '2026-02-04 13:45:08', '2026-02-04 13:45:03', '2026-02-04 13:45:08'),
(118, 8, 97, 'node_1', '新节点', 'inbound_order', 15, 'RK20260204135013', 3, '鞠磊', 'pending', false, '2026-02-04 13:50:13', NULL, '2026-02-04 13:50:13', '2026-02-04 13:50:13'),
(119, 8, 97, 'node_1', '新节点', 'inbound_order', 15, 'RK20260204135013', 4, '李波', 'pending', false, '2026-02-04 13:50:13', NULL, '2026-02-04 13:50:13', '2026-02-04 13:50:13'),
(117, 8, 97, 'node_1', '新节点', 'inbound_order', 15, 'RK20260204135013', 1, '系统管理员', 'approved', false, '2026-02-04 13:50:13', '2026-02-04 13:50:16', '2026-02-04 13:50:13', '2026-02-04 13:50:16');

-- Re-enabling triggers and foreign keys
SET session_replication_role = 'origin';
