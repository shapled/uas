create table uas_app (
  id bigint not null auto_increment primary key comment '主键 id',
  app varchar(16) not null comment '应用名称',
  description varchar(1024) comment '应用描述',
  status int default 0 comment '状态，0 启用，1 禁用，2 删除',
  created_by bigint not null comment '创建者 id',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  updated_at datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
  deleted_at datetime comment '删除时间'
);

create table uas_role (
  id bigint not null auto_increment primary key comment '主键 id',
  app_id bigint not null comment '关联的应用 id',
  role varchar(16) not null comment '角色名称',
  description varchar(1024) comment '角色描述',
  status int default 0 comment '状态，0 启用，1 禁用，2 删除',
  created_by bigint not null comment '创建者 id',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  updated_at datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
  deleted_at datetime comment '删除时间'
);

create table uas_permission (
  id bigint not null auto_increment primary key comment '主键 id',
  permission varchar(16) not null comment '权限名称',
  description varchar(1024) comment '权限描述',
  status int default 0 comment '状态，0 启用，1 禁用，2 删除',
  created_by bigint not null comment '创建者 id',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  updated_at datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
  deleted_at datetime comment '删除时间'
);

create table uas_role_permission (
  id bigint not null auto_increment primary key comment '主键 id',
  role_id bigint not null comment '角色 id',
  permission_id bigint not null comment '权限 id',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  unique key `uniq_user_role` (`role_id`, `permission_id`)
);

create table uas_user (
  id bigint not null auto_increment primary key comment '主键 id',
  nickname varchar(16) not null comment '用户昵称',
  username varchar(64) not null comment '登陆名称',
  password varchar(128) not null comment '用户密码',
  phone varchar(16) comment '联系电话',
  email varchar(64) comment '电子邮箱',
  extra json comment '其他字段',
  created_by bigint not null comment '创建者应用的 app id',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  updated_at datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
  deleted_at datetime comment '删除时间',
  unique key `uniq_phone_email` (`phone`, `email`)
);

create table uas_user_app (
  id bigint not null auto_increment primary key comment '主键 id',
  user_id bigint not null comment '用户 id',
  app_id bigint not null comment '应用 id',
  expired_at datetime comment '用户过期时间',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  updated_at datetime default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新时间',
  unique key `uniq_user_app` (`user_id`, `app_id`)
);

create table uas_user_app_role (
  id bigint not null auto_increment primary key comment '主键 id',
  user_id bigint not null comment '用户 id',
  app_id bigint not null comment '应用 id',
  role_id bigint not null comment '角色 id',
  created_at datetime default CURRENT_TIMESTAMP comment '创建时间',
  unique key `uniq_user_role` (`user_id`, `app_id`, `role_id`)
);
