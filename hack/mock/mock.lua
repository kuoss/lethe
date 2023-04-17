function node(tag, timestamp, r)
    local hostnames = {'node01', 'node02'}
    local syslog_identifiers = {'kubelet', 'containerd'}
    local messages = {'hello', 'world'}

    local u = os.date("%Y-%m-%dT%H:%M:%SZ")
    local hostname = hostnames[math.random(#hostnames)]
    local syslog_identifier = syslog_identifiers[math.random(#syslog_identifiers)]
    local message = messages[math.random(#messages)] .. " rand_value=" .. r.rand_value

    return 1, timestamp, {
        tag = string.format("/node/%s/%s_%s.log", hostname, string.sub(u, 1, 10), string.sub(u, 12, 13)),
        row = string.format("%s[%s|%s] %s", u, hostname, syslog_identifier, message)
    }
end

function pod(tag, timestamp, r)
    local namespaces = {'default', 'kube-system'}
    local pods = {'pod1', 'pod2'}
    local containers = {'main', 'sidecar'}
    local files = {'hello.go', 'world.go'}
    local keys = {'foo', 'bar'}

    local u = os.date("%Y-%m-%dT%H:%M:%SZ")
    local namespace = namespaces[math.random(#namespaces)]
    local pod = pods[math.random(#pods)]
    local container = containers[math.random(#containers)]

    local file = files[math.random(#files)]
    local key = keys[math.random(#keys)]
    local log = file .. ' ' .. key .. '=' .. r.rand_value

    return 1, timestamp, {
        tag = string.format("/pod/%s/%s_%s.log", namespace, string.sub(u, 1, 10), string.sub(u, 12, 13)),
        row = string.format("%s[%s|%s|%s] %s", u, namespace, pod, container, log)
    }
end
