# Features
1. monkey-patch provides the following features to patch golang functions and instance methods.
    - **AddPatchFunc** provides the functionality to patch any golang function
    - **AddPatchInstanceMethod** provides the functionality to patch any golang instance method

2. monkey-patch provides modification of patch behaviours via a configuration file (data.json).
    - modify patch values via data.json
    - enabling/disabling patches via data.json

# Installation
Get the latest monkey-patch in your repo
```sh
go get github.com/weilianlow/monkey-patch
```

# Create Patch Code
Patch code should be created at the package level to avoid unnecessary code refactoring due to unexported identifiers.
1. The following code snippets provide sample usage of **AddPatchFunc** and **AddPatchInstanceMethod**.
    - AddPatchFunc
      ```sh
      import mp "github.com/weilianlow/monkey-patch"

      func init(){
        mp.AddPatchFunc("ServiceGetBatchItem", mp.PatchMeta{
          Target: GetBatchItem,
          Replacement: func(data *mp.Data) interface{} {
            return func(ctx context.Context, shopItemIDs []*ShopItemID, needDeleted, needCensoringData bool) (_ map[int64]*bass_item.Item, _ int32) {
              // call the original function
              data.Guard.Unpatch()
              defer data.Guard.Restore()
              // get item basics from rpc request
              itemBasics, code := GetBatchItem(ctx, shopItemIDs, needDeleted, needCensoringData)
              // patch video infos into sorted item basic
              videoInfos := make([]*bass_item.VideoInfo, 0)
              json.Unmarshal([]byte(data.Value.(string)), &videoInfos)
              keys := make([]int, 0)
              for key, _ := range itemBasics {
                keys = append(keys, int(key))
              }
              sort.Ints(keys)
              for i, key := range keys {
                if len(videoInfos) > i {
                  itemBasics[int64(key)].HighlightVideo = videoInfos[i]
                }
              }
              return itemBasics, code
            }
          },
        })
      }
      ```

    - AddPatchInstanceMethod
      ```sh
      import mp "github.com/weilianlow/monkey-patch"

      func init(){
        mp.AddPatchInstanceMethod("GetBatchDispatcherInfo", mp.PatchMeta{
          Target:     &videoDispatcherServiceRPCImpl{},
          MethodName: "GetBatchDispatcherInfo",
          Replacement: func(data *mp.Data) interface{} {
            return func(_ *videoDispatcherServiceRPCImpl, _ context.Context, _ []byte) (_ []byte, _ uint32, _ error) {
              if val, ok := data.Value.([]interface{}); ok {
                code, _ := strconv.Atoi(val[1].(string))
                var err error
                if len(val[2].(string)) > 0 {
                  err = fmt.Errorf("GetBatchDispatcherInfo error '%v'", val[2])
                }
                return []byte(val[0].(string)), uint32(code), err
              }
              return nil, uint32(0), nil
            }
          },
        })
      }
      ```
2. Run your server. You should expect the following logs.
    ```sh
    2021-10-10 21:13:09.635846|DEBUG|-|96614:0:1|-|patch.go:73|-:monkey-patch.AddPatchInstanceMethod|-|successfully added data{name: 'GetBatchDispatcherInfo'} as patch method|-:-|
    2021-10-10 21:13:09.652621|DEBUG|-|96614:0:1|-|patch.go:49|-:monkey-patch.AddPatchFunc|-|successfully added data{name: 'ServiceGetBatchItem'} as patch function|-:-|
    ```

# Configure patch behaviours with config file
Patch behaviours are controlled by a config file (data.json). You can create the file in your root folder, etc folder or patch folder.
1. The following is just a sample based on the added patch function and instance method earlier on in this README file.
    ```sh
    {
      "data": [
        {
          "name": "ServiceGetBatchItem",
          "value": "[{\"video_id\":\"1\",\"thumb_url\":\"thumb_url\",\"duration\":100000,\"version\":2,\"vid\":\"one.mp4\",\"formats\":[{\"format\":\"mp4\",\"url\":\"format1\",\"width\":1280,\"height\":720},{\"format\":\"mov\",\"url\":\"format1\",\"width\":1280,\"height\":720}],\"default_format\":{\"format\":\"mp4\",\"url\":\"format1\",\"width\":1280,\"height\":720}}]",
          "enabled": true
        },
        {
          "name": "GetBatchDispatcherInfo",
          "value": [
            "{\"code\":0,\"msg\":\"success\",\"data\":[{\"vid\":\"sg_0077b29f-1e09-4143-9cf0-c649e3af1c69_000000\",\"duration\":0,\"formats\":[{\"format\":460011,\"defn\":\"V720P\",\"profile\":\"MP4\",\"path\":\"10520122104/sg_0077b29f-1e09-4143-9cf0-c649e3af1c69_000000.60011.mp4\",\"width\":1280,\"height\":720}],\"default_format\":{\"format\":200001,\"defn\":\"ORI\",\"profile\":\"MP4\",\"url\":\"\",\"width\":1920,\"height\":1080}}]}",
            "200",
            ""
          ],
          "enabled": true
        }
      ]
    }
    ```

2. Add the following code in your request handler, or middleware request function
    ```sh
    import mp "github.com/weilianlow/monkey-patch"

    func (cpm *CommonParametersMiddleware) ProcessRequest(pctx *middleware.ProcessorContext) middleware.ResultStatus {
    	mp.New(pctx.Ctx).MonkeyPatchByConfig(pctx.Ctx)
        ...
    }
    ```

3. That's it! You may try to trigger a request to check if the patch works.
    ```sh
    2021-10-10 21:13:28.375201|DEBUG|8b168866cdff60d325b208f7bdf8ac00:000000bfd9e2f424:0000000000000000|96614:0:204|-|patch.go:98|-:monkey-patch.(*DataList).monkeyPatch|-|enabling patch function for data{name: 'ServiceGetBatchItem'}|-:-|idc=sg2,func=monkeyPatch,spanID=010000de2169e6d7
    2021-10-10 21:13:28.375263|DEBUG|8b168866cdff60d325b208f7bdf8ac00:000000bfd9e2f424:0000000000000000|96614:0:204|-|patch.go:100|-:monkey-patch.(*DataList).monkeyPatch|-|enabling patch method for data{name: 'GetBatchDispatcherInfo'}|-:-|idc=sg2,func=monkeyPatch,spanID=010000de2169e6d7
    ```
