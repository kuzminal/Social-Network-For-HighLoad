box.cfg { listen = 3301 }

s = box.schema.space.create('posts', { if_not_exists = true })
s:format({
    { name = 'user_id', type = 'string' },
    { name = 'posts', type = 'array' }
})
s:create_index('primary', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'user_id' }
})

s = box.schema.space.create('users', { if_not_exists = true })
s:format({
    { name = 'Id', type = 'string' },
    { name = 'FirstName', type = 'string' },
    { name = 'SecondName', type = 'string' },
    { name = 'Age', type = 'unsigned' },
    { name = 'Birthdate', type = 'string' },
    { name = 'Biography', type = 'string' },
    { name = 'City', type = 'string' },
    { name = 'Password', type = 'string' },
})
s:create_index('primary', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'Id' }
})

s:create_index('search', {
    if_not_exists = true,
    type = 'TREE',
    unique = false,
    parts = {
        { 'FirstName' },
        { 'SecondName' }
    }
})

s:create_index('search', {
    if_not_exists = true,
    type = 'TREE',
    unique = false,
    parts = {
        {field = 2, type = 'string'},
        {field = 3, type = 'string'}
    }
})

s = box.schema.space.create('sessions', { if_not_exists = true })
s:format({
    { name = 'id', type = 'string' },
    { name = 'user_id', type = 'string' },
    { name = 'token', type = 'string' },
    { name = 'created_at', type = 'string' }
})
s:create_index('primary', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'id' }
})
s:create_index('token', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'token' }
})
s:create_index('user', {
    if_not_exists = true,
    type = 'hash',
    parts = { 'user_id' }
})


pg = require('pg')

require 'console'.start()
function update_cache_from_db(user_id)
    local conn = pg.connect({
        host = 'db',
        port = 5432,
        user = 'postgres',
        password = 'postgres',
        db = 'postgres'
    })
    local test = conn:execute("SELECT p.id, p.text, p.author_user_id, p.created_at " ..
            "FROM social.friends f " ..
            "RIGHT JOIN social.posts p ON p.author_user_id=f.friend_id " ..
            "WHERE f.user_id='" .. user_id .. '\' ' ..
            "ORDER BY p.created_at DESC LIMIT 1000; ")
    local row = ''
    for _, card in ipairs(test) do
        --box.space.posts:replace { user_id, card }
        --row = card
        row = box.space.posts:replace { user_id, card }
    end
    conn:close()
    return row
end

function get_data(key, offset, limit)
    local tuple = box.space.posts:get(key)
    local res = {}
    local response = {}
    if tuple == nil then
        tuple = update_cache_from_db(key)
        --tuple = box.space.posts:insert{key, data}
    end
    if tuple[2] ~= nil then
        res = tuple[2]
        if limit + offset > #res then
            limit = #res
        end
        for i = offset + 1, offset + limit do
            table.insert(response, res[i])
        end
    end
    return response
end

function create_user(id, firstName, secondName, age, birthdate, biography, city, password)
    res = box.space.users:insert {id, firstName, secondName, age, birthdate, biography, city, password}
    if res ~= nil then
        return res[1]
    end
    return ""
end

function get_user_by_id(key)
    local tuple = box.space.users:get(key)
    return tuple
end

function search_user(fnamepref, snamepref)
    local result = {}
    for _, fname in box.space.users.index.search:pairs(fnamepref, { iterator = 'GE', after = after}) do
        if (string.sub(fname[2], 1, string.len(fnamepref)) == fnamepref) and (string.sub(fname[3], 1, string.len(snamepref)) == snamepref) then
            table.insert(result, fname)
        else
            break
        end
    end
    return result
end

function check_user_exists(user)
    user_res = box.space.users:get(user)
    if user_res ~= nil then
        return true
    end
    return false
end

function get_session_by_user_id(token)
    user_res = box.space.sessions.index.token:get(token)
    if user_res ~= nil then
        return user_res[2]
    end
    return ""
end

function create_session(id, userId, token, created_at)
    tokenDb = box.space.sessions.index.user:get(userId)
    if tokenDb ~= nil then
        tok = tokenDb[3]
        return tokenDb[3]
    end
    res = box.space.sessions:insert { id, userId, token, created_at }
    if res ~= nil then
        return res[3]
    end
    return ""
end