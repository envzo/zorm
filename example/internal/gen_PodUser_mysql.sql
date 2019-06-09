use pod;

create table if not exists `pod_user` (
  `id` bigint not null auto_increment,
  `nickname` varchar(16) not null comment '昵称，非姓名',
  `password` varchar(16) not null,
  `age` int not null,
  `mobile_phone` varchar(18) not null default '13977q' comment '手机号',
  `sequence` bigint not null default '1' comment '顺序',
  `create_dt` bigint not null,
  `is_blocked` tinyint(1) not null,
  `update_dt` bigint not null comment '更新时间',
  `stats_dt` date not null,
  `dt` datetime not null default now(),
  primary key (`id`),
  unique key `uni_nickname_mobile_phone` (`nickname`, `mobile_phone`),
  unique key `uni_mobile_phone` (`mobile_phone`),
  index `idx_create_dt` (`create_dt`),
  index `idx_update_dt` (`update_dt`)
) engine=InnoDB default charset=utf8mb4 comment '我是无辜的测试表';
