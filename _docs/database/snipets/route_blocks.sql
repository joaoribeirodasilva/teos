INSERT INTO `app_routes_blocks` (`app_app_id`,`name`,`description`,`route`,`active`,`created_by`,`created_at`,`updated_by`,`updated_at`) VALUES
(2,'Autenticação','Autenticação de usuários','/auth',1,1,current_timestamp(),1,current_timestamp()),
(1,'Applicações','Aplicações disponíveis','/apps/apps',1,1,current_timestamp(),1,current_timestamp()),
(1,'Configurações','Configuração das aplicações','/apps/configurations',1,1,current_timestamp(),1,current_timestamp()),
(1,'Rotas','Rotas das aplicações','/apps/routes',1,1,current_timestamp(),1,current_timestamp()),
(1,'Blocos','Blocos das aplicações','/apps/routesblocks',1,1,current_timestamp(),1,current_timestamp()),
(4,'Histórico','Histórico de alterações de dados','/hists',1,1,current_timestamp(),1,current_timestamp()),
(4,'Histórico','Histórico de alterações de dados','/hists',1,1,current_timestamp(),1,current_timestamp()),
(3,'Usuários','Cadastros de usuários','/users/users',1,1,current_timestamp(),1,current_timestamp()),
(3,'Grupos','Grupos de usuários de usuários','/users/groups',1,1,current_timestamp(),1,current_timestamp()),
(3,'Perfis','Perfis de usuários de usuários','/users/roles',1,1,current_timestamp(),1,current_timestamp()),
(3,'Perfis de grupo','Perfis de grupos de usuários de usuários','/users/rolegroups',1,1,current_timestamp(),1,current_timestamp()),
(3,'Tipos de resets','Tipos de resets de dados','/users/resettypes',1,1,current_timestamp(),1,current_timestamp()),
(3,'Resets','Resets de dados','/users/resets',1,1,current_timestamp(),1,current_timestamp()),
(3,'Sessões','Sessões dos usuários','/users/sessions',1,1,current_timestamp(),1,current_timestamp())

INSERT INTO `app_routes` (`app_routes_block_id`,`name`,`description`,`method`,`open`,`active`,`created_by`,`created_at`,`updated_by`,`updated_at`) VALUES 
(1,`Login`,`Entrada no sistema, autenticação de usuário`,`POST`,`/login`,1,1,1,current_timestamp(),1,current_timestamp()),
(1,`Pedido de reset`,`Pedido de reset de senha pelo usuário`,`POST`,`/forgot`,1,1,1,current_timestamp(),1,current_timestamp()),
(1,`Reset de senha`,`Efetua o reset da senha do usuário`,`POST`,`/reset`,1,1,1,current_timestamp(),1,current_timestamp()),
(1,`Logout`,`Saída do sistema`,`DELETE`,`/logout`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Listagem`,`Listagem das aplicações do sistema`,`GET`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Cadastro`,`Cadastro de aplicação do sistema`,`GET`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Criar`,`Cria aplicação do sistema`,`POST`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Atualizar`,`Atualiza aplicação do sistema`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Apagar`,`Apaga aplicação do sistema`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),

(1,`Listagem`,`Listagem das configurações das aplicações do sistema`,`GET`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Cadastro`,`Cadastro da configuração das aplicação do sistema`,`GET`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Criar`,`Cria configuração de aplicação do sistema`,`POST`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Atualizar`,`Atualiza configuração de aplicação do sistema`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Apagar`,`Apaga configuração de aplicação do sistema`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),

(1,`Listagem`,`Listagem das rotas das aplicações`,`GET`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Cadastro`,`Rota de aplicação`,`GET`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Criar`,`Cria uma rota de aplicação`,`POST`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Atualizar`,`Atualiza uma rota de aplicação`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Apagar`,`Apaga uma rota de aplicação`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),

(1,`Listagem`,`Listagem blocos de rotas das aplicações`,`GET`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Cadastro`,`Bloco de rotas de aplicação`,`GET`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Criar`,`Cria um bloco de rotas de aplicação`,`POST`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Atualizar`,`Atualiza um bloco de rotas de aplicação`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(1,`Apagar`,`Apaga um bloco de rotas de aplicação`,`PUT`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),

(4,`Listagem`,`Listagem de alterações de dados`,`GET`,``,0,1,1,current_timestamp(),1,current_timestamp()),
(4,`Cadastro`,`Alteração de dados`,`GET`,`:id`,0,1,1,current_timestamp(),1,current_timestamp()),
(4,`Criar`,`Cria uma alteração de dados`,`POST`,``,0,1,1,current_timestamp(),1,current_timestamp()),




