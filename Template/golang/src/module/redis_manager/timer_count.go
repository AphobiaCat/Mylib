package redis_manager

import(
	"mylib/src/public"
    "time"
)

/*
	ret >= 0 is true
*/

func (rm *Redis_Manager) Timer_Count(redis_key string, reload_count int64, reset_time_s int64)(int64){
	ok, err := rm.rdb.SetNX(rm.ctx, redis_key, reload_count - 1, time.Duration(reset_time_s * 1000 * 1000 * 1000)).Result()
    if err != nil {
    	public.DBG_ERR("Timer_Count err[", err, "]")
        return -1
    }

    if ok {
        // first call init succ.
        return reload_count - 1
    }

    count, err := rm.rdb.Decr(rm.ctx, redis_key).Result()
    if err != nil {
        public.DBG_ERR("Timer_Count err[", err, "]")
        return -1
    }

    if count < 0 {
        // count done
        ttl, err := rm.rdb.TTL(rm.ctx, redis_key).Result()
	    if err == nil && ttl < 0 {
	        rm.rdb.Del(rm.ctx, redis_key)
	        return reload_count
	    }
        return -1
    }

    return count
}

func Timer_Count(redis_key string, reload_count int64, reset_time int64)(int64){
	return redis_manager.Timer_Count(redis_key, reload_count, reset_time)
}
