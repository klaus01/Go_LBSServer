var gBeginDateTime = null;
var gEndDateTime = null;
var gPayStatus = null;
var gDeliveryStatus = null;
var gPageNo = 1;
var gPageTotal = 1;
var gDatas = null;
const gPageSize = 20;
const gDateTimeFormat = "YYYY-MM-DD HH:mm:ss";

function reloadDatas(
  pageNo,
  pageSize,
  beginDateTime,
  endDateTime,
  payStatus,
  deliveryInfo
) {
  var searchButton = $("#searchButton").button("loading");
  $.get("/orders", {
    pageNo: pageNo,
    pageSize: pageSize,
    beginDateTime: beginDateTime,
    endDateTime: endDateTime,
    payStatus: payStatus,
    isShipped: deliveryInfo
  })
    .done(function(data) {
      if (data.code != 200) {
        var msg = data.message;
        if (msg == undefined || msg == null || msg.length <= 0) {
          alert("错误 " + data.code);
        } else {
          alert(msg);
        }
        return;
      }
      gBeginDateTime = beginDateTime;
      gEndDateTime = endDateTime;
      gPayStatus = payStatus;
      gDeliveryStatus = deliveryInfo;
      gPageNo = data.data.pageNo;
      gPageTotal = data.data.pageTotal;
      gDatas = data.data.datas;
      showData(data.data);
      shipButtonBindEvent();
      showPageNavigation(data.data);
      pageNavigationBindEvent();
    })
    .fail(function(res) {
      alert(
        `${res.status} ${res.statusText} ${
          res.responseJSON ? res.responseJSON.message : ""
        }`
      );
    })
    .always(function() {
      searchButton.button("reset");
    });
}

function showData(data) {
  var bodyHtml = "";
  var orders = data.datas;
  for (i in orders) {
    var order = orders[i];
    var createAt = moment(order.createAt).format(gDateTimeFormat);
    var payInfo = "";
    var payClass = "";
    if (order.payStatus == 0) {
      payInfo = "待支付";
    } else {
      payInfo =
        `${order.payInfo.modeName}|${order.payInfo.money}|` +
        moment(order.payInfo.createAt).format(gDateTimeFormat);
      payClass = order.payStatus == 1 ? "bg-success" : "bg-danger";
    }

    var deliveryInfo = "";
    if (order.deliveryInfo == undefined || order.deliveryInfo == null) {
      deliveryInfo = `未发货 <button name="shipButton" type="button" class="btn btn-primary btn-xs" dataIndex="${i}">发货</button>`;
    } else {
      deliveryInfo =
        `${order.deliveryInfo.courierCompany}|${order.deliveryInfo.waybillNumber}|` +
        moment(order.deliveryInfo.createAt).format(gDateTimeFormat);
    }

    bodyHtml += "<tr>";
    bodyHtml += `<td>${order.orderId}</td>`;
    bodyHtml += `<td>${createAt}</td>`;
    bodyHtml += `<td>${order.productName}</td>`;
    bodyHtml += `<td style="text-align:right">${order.productAmount}</td>`;
    bodyHtml += `<td>${order.name}</td>`;
    bodyHtml += `<td>${order.phoneNumber}</td>`;
    bodyHtml += `<td>${order.province} ${order.city} ${
      order.district == null ? "" : order.district
    } ${order.address}</td>`;
    bodyHtml += `<td class="${payClass}">${payInfo}</td>`;
    bodyHtml += `<td>${deliveryInfo}</td>`;
    bodyHtml += "</tr>";
  }
  $("table tbody").html(bodyHtml);
}

function shipButtonBindEvent() {
  $("button[name=shipButton]").on("click", function() {
    var dataIndex = parseInt($(this).attr("dataIndex"));
    if (dataIndex < 0 || dataIndex >= gDatas.length) return;
    var order = gDatas[dataIndex];
    var info = `${order.name} ${order.phoneNumber} ${order.orderId}`;
    var orderInfoNode = $("#shipDialog-orderInfo");
    orderInfoNode.text(info);
    orderInfoNode.attr("data-orderId", order.orderId);
    $("#shipDialog").modal("show");
  });
}

function showPageNavigation(data) {
  var pageNo = data.pageNo;
  var pageTotal = Math.max(data.pageNo, data.pageTotal);
  var leftClass = pageNo <= 1 ? "disabled" : "";
  var leftData = pageNo - 1;
  var rightClass = pageNo >= pageTotal ? "disabled" : "";
  var rightData = pageNo + 1;
  var bodyHtml = "";
  bodyHtml += `<li class="${leftClass}" data="${leftData}"><a href="#" aria-label="Previous"><span aria-hidden="true">&laquo;</span></a></li>`;
  for (var i = 1; i <= pageTotal; i++) {
    var classStr = i == pageNo ? "active" : "";
    bodyHtml += `<li class="${classStr}" data="${i}"><a href="#">${i}</a></li>`;
  }
  bodyHtml += `<li class="${rightClass}" data="${rightData}"><a href="#" aria-label="Next"><span aria-hidden="true">&raquo;</span></a></li>`;
  $("#pageNavigation").html(bodyHtml);
}

function pageNavigationBindEvent() {
  $("#pageNavigation li").on("click", function() {
    var newPageNo = parseInt($(this).attr("data"));
    if (newPageNo == gPageNo) return;
    if (newPageNo < 1) return;
    if (newPageNo > gPageTotal) return;
    reloadDatas(
      newPageNo,
      gPageSize,
      gBeginDateTime,
      gEndDateTime,
      gPayStatus,
      gDeliveryStatus
    );
  });
}

$(function() {
  $("#beginDateTime").datetimepicker({
    format: "YYYY-MM-DD",
    locale: moment.locale("zh-cn")
  });
  $("#endDateTime").datetimepicker({
    format: "YYYY-MM-DD",
    locale: moment.locale("zh-cn")
  });
  $("#payStatusList li").on("click", function() {
    var statusText = $(this).text();
    var status = $(this).attr("data");
    var button = $("#payStatusButton");
    status == undefined
      ? button.removeAttr("data")
      : button.attr("data", status);
    button.html("支付状态：" + statusText + ' <span class="caret"></span>');
  });

  $("#deliveryStatusList li").on("click", function() {
    var statusText = $(this).text();
    var status = $(this).attr("data");
    var button = $("#deliveryStatusButton");
    status == undefined
      ? button.removeAttr("data")
      : button.attr("data", status);
    button.html("发货状态：" + statusText + ' <span class="caret"></span>');
  });

  $("#searchButton").on("click", function() {
    var beginDateTime = $("#beginDateTime").val();
    if (beginDateTime.length <= 0) beginDateTime = undefined;
    else beginDateTime = moment(beginDateTime).format();
    var endDateTime = $("#endDateTime").val();
    if (endDateTime.length <= 0) endDateTime = undefined;
    else
      endDateTime = moment(endDateTime)
        .add(1, "days")
        .format();
    var payStatus = $("#payStatusButton").attr("data");
    var deliveryStatus = $("#deliveryStatusButton").attr("data");

    reloadDatas(
      1,
      gPageSize,
      beginDateTime,
      endDateTime,
      payStatus,
      deliveryStatus
    );
  });
  $("#shipSubmit").on("click", function() {
    var orderId = $("#shipDialog-orderInfo").attr("data-orderId");
    if (orderId == undefined) {
      alert("缺少订单ID");
      return;
    }
    var $btn = $(this).button("loading");
    $.post(
      "/orders/deliveryInfo",
      {
        orderId: orderId,
        courierCompany: $("#shipInputCourierCompany").val(),
        waybillNumber: $("#shipInputWaybillNumber").val()
      },
      "json"
    )
      .done(function(data) {
        if (data.code != 200) {
          var msg = data.message;
          if (msg == undefined || msg == null || msg.length <= 0) {
            alert("错误 " + data.code);
          } else {
            alert(msg);
          }
        } else {
          $("#shipDialog").modal("hide");
          reloadDatas(
            gPageNo,
            gPageSize,
            gBeginDateTime,
            gEndDateTime,
            gPayStatus,
            gDeliveryStatus
          );
        }
      })
      .fail(function(res) {
        alert(
          `${res.status} ${res.statusText} ${
            res.responseJSON ? res.responseJSON.message : ""
          }`
        );
      })
      .always(function() {
        $btn.button("reset");
      });
  });
});
