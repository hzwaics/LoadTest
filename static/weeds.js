var my_chart;
var my_chart_server;
var weeds_status;//0:stop 1:start;
var chart_controller;
var req_num = []
var succ_resp_num = []
var delay = []
var fail_num = []
var date_string = []
var cpu = []
var memory = []
var server_date_string = []

function get_test_data(){
	$.get("get_second_data",function(data,status){
		req_num.push(JSON.parse(data)['one_second_data']['request_num'])
		succ_resp_num.push(JSON.parse(data)['one_second_data']['success_response_num'])
		delay.push(JSON.parse(data)['one_second_data']['average_cost_time'])
		fail_num.push(JSON.parse(data)['one_second_data']['failed_num'])
		time_stamp = JSON.parse(data)['one_second_data']['time_stamp']
		var date = new Date(time_stamp*1000);
		Y = date.getFullYear() + '-';
		M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
		D = date.getDate() + ' ';
		h = date.getHours() + ':';
		m = date.getMinutes() + ':';
		s = date.getSeconds();
		date_string.push(Y+M+D+h+m+s)
		opt = {
			xAxis:{
				data:date_string
			},
			series:[{
				name:'请求数',
				data:req_num
			},
				{
					name:'成功响应数',
					data:succ_resp_num
				},
				{
					name:'失败数',
					data:fail_num
				},
				{
					name:'响应延时',
					yAxisIndex:1,
					data:delay
				},
			]
		}
		//console.log(opt);
		my_chart.setOption(opt);
		if(JSON.parse(data)['target_server_data']['time_stamp'] != 0){
			cpu.push(JSON.parse(data)['target_server_data']['cpu_rate']/100.0)
			memory.push(JSON.parse(data)['target_server_data']['memory_usage']/100.0)
			server_time_stamp = JSON.parse(data)['target_server_data']['time_stamp']
			var date = new Date(server_time_stamp*1000);
			Y = date.getFullYear() + '-';
			M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
			D = date.getDate() + ' ';
			h = date.getHours() + ':';
			m = date.getMinutes() + ':';
			s = date.getSeconds();
			server_date_string.push(Y+M+D+h+m+s)
			opt = {
				xAxis:{
					data:server_date_string
				},
				series:[{
					name:'CPU占用',
					data:cpu
				},
					{
						name:'内存占用',
						data:memory
					},
				]
			}
			my_chart_server.setOption(opt);
		}
		$("#req_num_sum").text(JSON.parse(data)['test_data_sum']['total_request_num'])
		$("#succ_resp_num_sum").text(JSON.parse(data)['test_data_sum']['total_succ_response_num'])
		$("#failed_num_sum").text(JSON.parse(data)['test_data_sum']['total_failed_num'])
		server_status = JSON.parse(data)['server_status']
		if (server_status == 0) {
			$("#status").text("Stop")
		}
		else if (server_status == 1){
			$("#status").text("Running: " + JSON.parse(data)['one_second_data']['request_num'] + " requests/sec")
		}
		slave_num = JSON.parse(data)['slave_num']
		if (slave_num > 0){
			$("#slave").text(slave_num +' Slaves');
		}
	});
}

function start_test(){
	if (weeds_status == 1) {
		alert("请先停下当前测试任务，再开始新的任务");
		return
	}
	$.post("start",JSON.stringify({"target-qps":Number($("#target-qps").val()),"increase-per-second":Number($("#increase-per-second").val())}),function (){
	
		weeds_status = 1;
		chart_controller = setInterval(get_test_data, 1000);
		show_chart();
	
});
}
function stop_test(){
	$.post("stop","stop",function (){
		setTimeout(function(){
			clearInterval(chart_controller);
			show_start();
			weeds_status = 0
		},1000);
	
});
}


function show_start(){
	$(".nav").css("background-color","")
	$("#start_nav").css("background-color","#e7e7e7")
	$("#start").show();
	$("#chart_div").hide();
	$("#sum_div").hide();
}

function show_chart(){
	$(".nav").css("background-color","")
	$("#chart_nav").css("background-color","#e7e7e7")
	$("#start").hide();
	$("#sum_div").hide();
	$("#chart_div").show();
}

function show_sum(){
	$(".nav").css("background-color","")
	$("#sum_nav").css("background-color","#e7e7e7")
	$("#start").hide();
	$("#sum_div").show();
	$("#chart_div").hide();
}

$(document).ready(function(){
	show_start();
	my_chart = echarts.init(document.getElementById('chart'));
	my_chart.setOption({
		title:{
			text:'请求及延时'
		},
		tooltip:{
			trigger: 'axis'
		},
		legend:{
			data:['请求数','成功响应数','响应延时','失败数']
		},
		xAxis:{
			boundaryGap: false,
			data:[]
		},
		yAxis:[
			{
				name: "请求/响应数",
				type: "value"
			},
			{
				name: "响应延时",
				type: "value"
			}
		],
		series:[{
			name:'请求数',
			type:'line',
			data:[]
			},
			{
			name:'成功响应数',
			type:'line',
			data:[]
			},
			{
			name:'响应延时',
			type:'line',
			yAxisIndex:1,
			data:[]
			},
			{
			name:'失败数',
			type:'line',
			data:[]
			},
		]
	});
	my_chart_server = echarts.init(document.getElementById('chart_server'));
	my_chart_server.setOption({
		title:{
			text:'目标服务器资源'
		},
		tooltip:{
			trigger: 'axis'
		},
		legend:{
			data:['CPU占用','内存占用']
		},
		xAxis:{
			boundaryGap: false,
			data:[]
		},
		yAxis:[
			{
				name: "CPU占用(100%-idle)",
				type: "value"
			},
			{
				name: "内存占用(100%-avaiable)",
				type: "value"
			}
		],
		series:[{
			name:'CPU占用',
			type:'line',
			data:[]
			},
			{
			name:'内存占用',
			type:'line',
			yAxisIndex:1,
			data:[]
			}
		]
	});
	$.get("get_second_data",function(data,status){
		server_status = JSON.parse(data)['server_status']
		if (server_status == 0) {
			$("#status").text("Stop")
		}
		else if (server_status == 1){
			$("#status").text("Running: " + JSON.parse(data)['one_second_data']['request_num'] + " requests/sec")
			chart_controller = setInterval(get_test_data, 1000);
			show_chart();
			weeds_status = 1;
		}
		slave_num = JSON.parse(data)['slave_num']
		if (slave_num > 0){
			$("#slave").text(slave_num +' Slaves');
		}
		if (JSON.parse(data)['target_server_data']['time_stamp'] != 0){
			$("#chart_server").show();
		}
		else{
			$("#chart_server").hide();
		}
	});
});
