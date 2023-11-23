package main

import "practice/internal/config"

func getRedisConf(conf *config.AppConf) string {
	return conf.Redis.Host + ":" + conf.Redis.Port
}

func getKafkaConf(conf *config.AppConf) string {
	return conf.Kafka.Host + ":" + conf.Kafka.Port
}

func getGRPCConf(conf *config.AppConf) string {
	return ":" + conf.GRPC.Port
}
