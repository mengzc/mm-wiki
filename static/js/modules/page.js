/**
 * Copyright (c) 2018 phachon@163.com
 */

var Page = {

    /**
     * ajax Save
     * @param element
     * @returns {boolean}
     */
    ajaxSave: function (element) {

        /**
         * 成功信息条
         * @param message
         * @param data
         */
        function successBox(message, data) {
            Layers.successMsg(message)
        }

        /**
         * 失败信息条
         * @param message
         * @param data
         */
        function failed(message, data) {
            Layers.failedMsg(message)
        }

        /**
         * request success
         * @param result
         */
        function response(result) {
            //console.log(result)
            if (result.code == 0) {
                failed(result.message, result.data);
            }
            if (result.code == 1) {
                successBox(result.message, result.data);
            }
            if (result.redirect.url) {
                var sleepTime = result.redirect.sleep || 3000;
                setTimeout(function () {
                    parent.location.href = result.redirect.url;
                }, sleepTime);
            }
        }

        var commentHtml =
            '<div class="container-fluid" style="padding: 20px 20px 0 20px">'+
                '<textarea name="edit_comment" class="form-control" rows="3" autofocus="autofocus" style="resize:none""></textarea>'+
                '<div style="margin-top: 8px;text-align: left">'+
                    '<label> 通知关注用户&nbsp;</label>' +
                    '<input type="checkbox" name="is_notice_user" checked="checked" value="1">' +
                '</div>' +
            '</div>';

        layer.open({
            title: '<i class="fa fa-volume-up"></i> 请输入修改备注',
            type: 1,
            area: ['380px', '232px'],
            content: commentHtml,
            btn: ['确定','取消'],
            yes: function(index, layero){
                var commentText = $("textarea[name='edit_comment']").val().trim();
                var isNoticeCheck = $("input[type='checkbox'][name='is_notice_user']").is(':checked');
                var isNoticeUser = "0";
                if (isNoticeCheck) {
                    isNoticeUser = "1";
                }
                if (commentText && commentText.length > 0) {
                    if (commentText.length > 50 ) {
                        layer.tips("最多50个字符！", $("textarea[name='edit_comment']"))
                    }else {
                        layer.close(index);
                        var options = {
                            dataType: 'json',
                            success: response,
                            data: {'comment': commentText, 'is_notice_user': isNoticeUser}
                        };
                        $(element).ajaxSubmit(options);
                    }
                }else {
                    $("textarea[name='edit_comment']").focus();
                }
            },
            btn2: function(index, layero){
                layer.close(index);
            }
        });

        // layer.prompt({
        //     title: '<i class="fa fa-volume-up"></i> 请输入修改备注',
        //     formType: 2,
        //     maxlength: 150,
        //     value: '',
        //     area: ['340px', '80px']
        // }, function(comment, index, elem){
        //     if (comment.trim()) {
        //         layer.close(index);
        //         var options = {
        //             dataType: 'json',
        //             success: response,
        //             data: {'comment': comment}
        //         };
        //         $(element).ajaxSubmit(options);
        //     }else {
        //         elem.focus()
        //     }
        // });

        return false;
    },

    /**
     * cancel save
     */
    cancelSave: function (title, url) {
        title = '<i class="fa fa-volume-up"></i> '+title;
        layer.confirm(title, {
            btn: ['是','否'],
            skin: Layers.skin,
            btnAlign: 'c',
            title: "<i class='fa fa-warning'></i><strong> 警告</strong>"
        }, function() {
            parent.location = url
        }, function() {

        });
    }

};