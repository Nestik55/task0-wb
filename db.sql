CREATE TABLE "order" (
    "order_uid" text PRIMARY KEY,
	"track_number" text UNIQUE,
	"entry" text NOT NULL,
	"locale" text NOT NULL,
	"internal_signature" text,
	"customer_id" text NOT NULL,
	"delivery_service" text NOT NULL,
	"shardkey" text NOT NULL,
	"sm_id" int NOT NULL,
	"date_created" text NOT NULL,
	"oof_shard" text NOT NULL
);

CREATE TABLE "delivery" (
    "order_id" text PRIMARY KEY REFERENCES "order"("order_uid"),
	"phone" text NOT NULL UNIQUE,
	"name" text NOT NULL,
	"zip" text NOT NULL,
	"city" text NOT NULL,
	"address" text NOT NULL,
	"region" text NOT NULL,
	"email" text NOT NULL
);

CREATE TABLE "payment" (
    "order_id" text PRIMARY KEY REFERENCES "order"("order_uid"),
	"transaction" text NOT NULL,
	"request_id" text,
	"currency" text NOT NULL,
	"provider" text NOT NULL,
	"amount" int NOT NULL,
	"payment_dt" int NOT NULL,
	"bank" text NOT NULL,
	"delivery_cost" int NOT NULL,
	"goods_total" int NOT NULL,
	"custom_fee" int NOT NULL
);

CREATE TABLE "item" (
	"id" SERIAL PRIMARY KEY,
    "order_id" text,
	"chrt_id" int NOT NULL,
	"track_number" text NOT NULL,
	"price" int NOT NULL,
	"rid" text NOT NULL,
	"name" text NOT NULL,
	"sale" int NOT NULL,
	"size" text NOT NULL,
	"total_price" int NOT NULL,
	"nm_id" int NOT NULL,
	"brand" text NOT NULL,
	"status" int NOT NULL
);