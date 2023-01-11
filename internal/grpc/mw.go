package grpc

// // constants ...
// const (
// 	// UserIDFieldName ...
// 	UserIDFieldName = "user"
// 	// RawUserIDFieldName ...
// 	RawUserIDFieldName = "raw_user"
// )
//
// type userIDCtxKey struct{}
//
// func AuthFromMD() grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (
// 		resp interface{},
// 		err error,
// 	) {
// 		var newCtx context.Context
// 		rawID := uuid.New().String()
//
// 		md, ok := metadata.FromIncomingContext(ctx)
// 		if ok {
// 			id := md.Get(UserIDFieldName)
// 			if len(id) == 0 {
// 				rawID = uuid.New().String()
// 			} else if len(id) > 1 {
// 				return nil, status.Error(codes.Unauthenticated, "got unexpected user id")
// 			} else if err := encryptor.Get().DecodeUUID(id[0], &rawID); err != nil {
// 				rawID = uuid.New().String()
// 			}
// 		}
// 		// md := metadata.New(map[string]string{
// 		// 	UserIDFieldName: rawID,
// 		// })
// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 		return handler(newCtx, req)
// 	}
// }
