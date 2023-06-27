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
