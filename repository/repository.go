package repository

import "sync"

type Repository interface {
	// ポイント等更新用のロック
	sync.Locker
	// ユーザーレーティング
	RatingRepository
	// スタンプエフェクトのポイント
	EffectPointRepository
	// スタンプ同士の案系
	StampRelationRepository
	// スタンプ単体のポイント
	StampRepository
	// 処理したチャンネルの保存用
	SeenChannelRepository
}
