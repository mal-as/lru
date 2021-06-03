box.cfg{listen = 3301}

local function init()
    box.schema.user.create('go', {
        password = 'passwd', 
        if_not_exists = true
    })

    box.schema.user.grant('go','create','universe',nil,{if_not_exists = true})
    box.schema.user.grant('go','write','universe',nil,{if_not_exists = true})
    box.schema.user.grant('go','read','universe',nil,{if_not_exists = true})
    box.schema.user.grant('go', 'execute', 'universe',nil,{if_not_exists = true})

    local s = box.schema.space.create('_lru_cache',{if_not_exists = true})

    s:format({
        {name = 'id', type = 'unsigned'},
        {name = 'key', type = 'string'},
        {name = 'data', type = 'string'}
    })
    box.schema.sequence.create('s', {min=0, start=0})
    s:create_index('primary', {sequence = 's'})
    s:create_index('secondary', {
        type = 'hash',
        parts = {'key'}
    })
end

-- глобальная переменная количества элементов в кэше, по умолчанию 100
CACHE_SIZE = 100

function setCacheSize(size)
    CACHE_SIZE = size
end

-- вставка данных в кэш
function set(key, value)
    -- пытаемся найти элемент с таким ключом
    local tuple = box.space._lru_cache.index.secondary:select{key}[1]

    --если элемент найден, то удаляем его
    if tuple ~= nil then
        delete(tuple)
    end

    -- вставляем в конец кортеж
    local tuples = box.space._lru_cache:insert{box.sequence.s:next(), key, value}

    -- если длинна спейса равна размеру кэша, то удаляем первый(самый старый) элемент
    if box.space._lru_cache:len() > CACHE_SIZE then
        local tuple = box.space._lru_cache:select()[1]
        delete(tuple)
    end

    return tuples
end

-- получение данных из кэша
function get(key)
    -- пытаемся найти элемент с таким ключом
    local tuple = box.space._lru_cache.index.secondary:select{key}[1]

    -- удаляем элемент и вставляем заново, что бы он оказался в начале кэша
    if tuple ~= nil then
        delete(tuple)
        tuple = set(tuple[2], tuple[3])
    end

    return tuple
end

-- удаление кортежа из спейса
function delete(tuple)
    -- выполняем удаление по первичному ключу
    box.space._lru_cache:delete(tuple[1])
end

function truncate()
    box.space._lru_cache:truncate()
end

box.once('init', init)
