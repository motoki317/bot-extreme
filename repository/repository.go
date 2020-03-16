package repository

type Repository interface {
	// ポイント等更新用のロック
	Lock()
	Unlock()
	// ユーザーレーティング
	RatingRepository
	// スタンプエフェクトのポイント
	EffectPointRepository
	// スタンプ同士の案系
	StampRelationRepository
	// スタンプ単体のポイント
	StampRepository
	// チャンネル更新用のロック
	ChannelLock()
	ChannelUnlock()
	// 処理したチャンネルの保存用
	SeenChannelRepository
}
