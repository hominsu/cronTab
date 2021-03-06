// 页面加载完成之后回调函数
$(document).ready(function () {
    // 1. 绑定按钮的事件处理函数
    // javascript 委托机制, DOM 事件冒泡的一个关键原理

    // 新建任务
    $("#new-job").on("click", newJobCallBack)

    // 健康节点
    $("#node").on("click", NodeCallBack)

    // 刷新按钮
    $("#refresh").on("click", refreshCallBack)

    const job_list = $("#job-list");

    // 编辑按钮
    job_list.on("click", ".edit-job", editJobCallBack)

    // 删除按钮
    job_list.on("click", ".delete-job", deleteJobCallBack)

    // 强杀按钮
    job_list.on("click", ".kill-job", killJobCallBack)

    // 日志按钮
    job_list.on("click", ".job-log", jobLogCallBack)

    // 编辑模态框中提交按钮
    $('#commit-job').on("click", commitJobCallBack)

    // 日志模态框中的删除按钮
    $('#delete-log').on("click", clearJobLogCallBack)

    // 刷新任务列表
    rebuildJobList()
})

// 新建任务按钮回调函数
function newJobCallBack() {
    // 清空模态框
    $('#edit-name').val("")
    $('#edit-command').val("")
    $('#edit-cron-expr').val("")

    // 弹出模态框
    $('#edit-modal').modal('show')
}

// 新建任务按钮回调函数
function NodeCallBack() {
    const $btn = $(this).button('loading');
    $.ajax({
        url: '/job/node',
        type: 'get',
        dataType: 'json',
        success: function (resp) {
            // 任务数组
            const nodeList = resp.data;

            // 清理列表
            const node_list_tbody = $('#node-list tbody')
            node_list_tbody.empty()

            if (nodeList !== null) {
                // 遍历任务, 填充 table
                for (let i = 0; i < nodeList.length; ++i) {
                    node_list_tbody.append($("<tr>").append($('<td>').html(nodeList[i])))
                }
            }

            // 弹出模态框
            $('#node-modal').modal('show')
        }
    })
    $btn.button('reset')
}

// 刷新按钮回调函数
function refreshCallBack() {
    const $btn = $(this).button('loading');
    // 刷新任务表格
    rebuildJobList()
    $btn.button('reset')

    const alert_success = $('#alert-success-modal #alert-success .alert-success-content');
    alert_success.empty()
    alert_success.append($("<p>").append("Msg: refresh success"))
    // 弹出模态框
    $('#alert-success-modal').modal('show')
}

// 修改按钮回调函数
function editJobCallBack() {
    // 取当前 job 的信息赋值给模态框的 input
    $('#edit-name').val($(this).parents('tr').children('.job-name').text())
    $('#edit-command').val($(this).parents('tr').children('.job-command').text())
    $('#edit-cron-expr').val($(this).parents('tr').children('.job-cron-expr').text())
    // 弹出模态框
    $('#edit-modal').modal('show')
}

// 删除按钮回调函数
function deleteJobCallBack() {
    const $btn = $(this).button('loading');
    const jobName = $(this).parents("tr").children(".job-name").text();
    $.ajax({
        url: '/job',
        type: 'delete',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify({name: jobName}),
        success: function (resp) {
            // 刷新任务表格
            rebuildJobList()

            const alert_success = $('#alert-success-modal #alert-success .alert-success-content');
            alert_success.empty()
            alert_success.append($("<p>").append("Msg: " + resp.msg))
            alert_success.append($("<p>").append("Delete Job: "))
            appendOldJob(alert_success, resp)
            // 弹出模态框
            $('#alert-success-modal').modal('show')
        }
    })
    $btn.button('reset')
}

// 强杀按钮回调函数
function killJobCallBack() {
    const $btn = $(this).button('loading');
    const jobName = $(this).parents("tr").children(".job-name").text();
    $.ajax({
        url: '/job/kill',
        type: 'post',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify({name: jobName}),
        success: function (resp) {
            const alert_success = $('#alert-success-modal #alert-success .alert-success-content');
            alert_success.empty()
            alert_success.append($("<p>").append("Msg: " + resp.msg))
            // 弹出模态框
            $('#alert-success-modal').modal('show')
        }
    })
    $btn.button('reset')
}

// 日志按钮回调函数
function jobLogCallBack() {
    const $btn = $(this).button('loading');
    const jobName = $(this).parents("tr").children(".job-name").text();
    $.ajax({
        url: '/job/log',
        type: 'post',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify({name: jobName}),
        success: function (resp) {
            // 任务数组
            const logList = resp.data;

            // 清理标题
            const title = $('#log-modal .modal-title')
            title.empty()
            title.html(jobName)

            // 清理列表
            const log_list_tbody = $('#log-list tbody')
            log_list_tbody.empty()

            // 日志列表不为空
            if (logList !== null) {
                // 遍历任务, 填充 table
                for (let i = 0; i < logList.length; ++i) {
                    const log = logList[i];

                    const tr = $("<tr>")
                    tr.append($('<td class="job-command">').html(log.command))
                    tr.append($('<td class="job-output">').html(log.output))
                    tr.append($('<td class="job-err">').html(log.err))
                    tr.append($('<td class="job-plan-time">').html(log.plan_time))
                    tr.append($('<td class="job-schedule-time">').html(log.schedule_time))
                    tr.append($('<td class="job-usage-time">').html(log.end_time - log.start_time))

                    log_list_tbody.append(tr)
                }
            }

            // 弹出模态框
            $('#log-modal').modal('show')
        }
    })
    $btn.button('reset')
}

// 日志模态框中清空按钮回调函数
function clearJobLogCallBack() {
    const $btn = $(this).button('loading');
    const jobName = $(this).parents('.modal-content').children('.modal-header').children('.modal-title').text();
    $.ajax({
        url: '/job/log',
        type: 'delete',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify({"name": jobName}),
        success: function (resp) {
            // 隐藏日志模态框
            $('#log-modal').modal('hide')

            if (resp.data !== 0) {
                const alert_success = $('#alert-success-modal #alert-success .alert-success-content');
                alert_success.empty()
                alert_success.append($("<p>").append("Msg: " + resp.msg))
                alert_success.append($("<p>").append("Del Job Count: " + resp.data))
                // 弹出成功模态框
                $('#alert-success-modal').modal('show')
            } else {
                const alert_danger = $('#alert-danger-modal #alert-danger .alert-danger-content')
                alert_danger.empty()
                alert_danger.append($("<p>").append("Msg: no job log to delete"))
                // 弹出错误模态框
                $('#alert-danger-modal').modal('show')
            }
        }
    })
    $btn.button('reset')
}

// 编辑模态框中提交按钮回调函数
function commitJobCallBack() {
    const $btn = $(this).button('loading');
    const jobInfo = {
        name: $('#edit-name').val(),
        command: $('#edit-command').val(),
        cron_expr: $('#edit-cron-expr').val()
    };
    $.ajax({
        url: '/job',
        type: 'post',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify(jobInfo),
        success: function (resp) {
            // 隐藏编辑模态框
            $('#edit-modal').modal('hide')

            // 刷新任务表格
            rebuildJobList()

            const alert_success = $('#alert-success-modal #alert-success .alert-success-content');
            alert_success.empty()
            alert_success.append($("<p>").append("Msg: " + resp.msg))
            if (resp.data !== null) {
                alert_success.append($("<p>").append("Prev Job: "))
                appendOldJob(alert_success, resp)
            }
            // 弹出模态框
            $('#alert-success-modal').modal('show')
        }
    })
    $btn.button('reset')
}

// 向成功模态警告框中加入上一次 job 的信息
function appendOldJob(alert_success, resp) {
    const p = "<p style='margin-left: 30px'>"
    alert_success.append($(p).append("job name: " + resp.data.name))
    alert_success.append($(p).append("command: " + resp.data.command))
    alert_success.append($(p).append("cron expr: " + resp.data.cron_expr))
}

// 刷新任务列表
function rebuildJobList() {
    $.ajax({
        url: '/job',
        type: 'get',
        dataType: 'json',
        success: function (resp) {
            // 服务端出错
            if (resp.errno !== 0) {
                const alert_danger = $('#alert-danger-modal #alert-danger .alert-danger-content')
                alert_danger.empty()
                alert_danger.append($("<p>").append("Msg: " + resp.msg))
                alert_danger.append($("<p>").append("Data: " + resp.data))
                // 弹出模态框
                $('#alert-danger-modal').modal('show')
                return
            }

            // 任务数组
            const jobList = resp.data;

            // 清理列表
            const job_list_tbody = $('#job-list tbody')
            job_list_tbody.empty()

            // 遍历任务, 填充 table
            for (let i = 0; i < jobList.length; ++i) {
                const job = jobList[i];

                const tr = $("<tr>")
                tr.append($('<td class="job-name">').html(job.name))
                tr.append($('<td class="job-command">').html(job.command))
                tr.append($('<td class="job-cron-expr">').html(job.cron_expr))

                const toolbar = $('<div class="btn-toolbar">')
                    .append($('<button class="btn btn-info edit-job">编辑</button>'))
                    .append($('<button class="btn btn-danger delete-job">删除</button>'))
                    .append($('<button class="btn btn-warning kill-job">强杀</button>'))
                    .append($('<button class="btn btn-success job-log">日志</button>'));

                tr.append($("<td>").append(toolbar))
                job_list_tbody.append(tr)
            }
        }
    })
}