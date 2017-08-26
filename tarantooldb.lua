#!/usr/bin/tarantool
-- Tarantool init script

local log = require('log')
local console = require('console')

box.cfg {
	listen = 3302,
}

if not box.space.examples then
	box.schema.space.create('examples',{id=999})
	box.space.examples:create_index('primary', {type = 'hash', parts = {1, 'NUM'}})
end

if not box.schema.user.exists("test") then
	box.schema.user.create('test', {password = '12345'})
	box.schema.user.grant('test','read,write','space','examples')
	box.schema.user.grant('test', 'read,write,execute', 'universe')
end