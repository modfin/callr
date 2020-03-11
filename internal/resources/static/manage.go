package static

import "github.com/labstack/echo"

func Manage(c echo.Context) error {
	return c.Blob(200, "text/html", []byte(manage))
}

const manage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Callr</title>
</head>
<link rel="stylesheet" href="//code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css">
<style>
	body{
		font-family: Arial, sans-serif ;
		padding: 0;
		margin: 0;
	}
	h1{
		text-align: center;
		margin: 0;
		padding: 0 0 25px 0;
		border-bottom: 2px solid #ccc;
    	font-size: 65px;
    	color: #777;
	}
	h1 img{
		height: 90px;
		margin: 14px 0 -16px;
	}
	a{
		cursor: pointer;
	}
	.container{
		margin: 10px;
		padding: 10px;
		border: 1px solid #ccc;
		min-height: 300px;
	}
	.history-item{
		margin: 10px;
		padding: 10px;
		border-top: 1px solid #CCC;
	}
	
	.links{
		margin-top: -10px;
		float: right;
		padding: 10px;
		text-align: right;
		line-height: 2em;
	}
	.links a{
		color: gray;
		text-decoration: none;
		border: 1px solid #ccc;
		border-radius: 2px;
		padding: 3px;
	}
	.links a:hover{
		color: #454545;
		text-decoration: none;
	}

	table{
		border-collapse: collapse;
	}
	table td, table th{
		text-align: left;
		padding: 2px 4px 2px 0;
	}
	form label{
		display: inline-block;
		width: 60px;
	}
	form input{
		display: inline-block;
		width: 300px;
	}
	button{
		cursor: pointer;
		background-color: #d5d5d5;
		color: #454545;
		font-weight: bold;
		border: 1px solid #7e7e7e;
		padding: 5px;
		box-shadow: 1px 1px 2px 1px rgba(0, 0, 0, 0.2)
	}
	button:hover{
		background-color: #e4e4e4;
		box-shadow: 1px 1px 3px 1px rgba(0, 0, 0, 0.2)
	}
	
	.saving-people, .saving-oncall {
		font-size: 13px;
		float: right;
		color: #6b6b6b;
		opacity: 0;
	}
	.saving-people.blink, .saving-oncall.blink{
		animation-name: blink;
  		animation-duration: 2s;
	}
	@keyframes blink {
	  0% {opacity: 0;}
 	  50%  {opacity: 1;}
	  100% {opacity: 0;}
	}
	#people-outlet{
		display: inline-block;
		width: 50%;
    	vertical-align: top;
	}
	#oncall-outlet .person{
		
	}
	#people-outlet .p-wrapper, #oncall-outlet .person{
		padding: 4px;
    	/*margin-bottom: -4px;*/
		background-color: #dbdbdb;
	}
	#people-outlet .p-wrapper:nth-child(even), #oncall-outlet .person:nth-child(even){
		background-color: #eaeaea;
	}
	#add-people{
		display: inline-block;
		width: 50%;
	}
	#add-people form{
		margin: 0 10px;
	}
	#add-people form p:first-child{
		margin-top: 0;
	}
	.p-links{
		font-size: 90%;
		color: #787878;
	}
	.p-links a:hover {
		color: #3a3a3a;
	}
</style>
<body>

<div>
	<div class="links">
		<a href="..">Incidents</a><br/>
	</div>
</div>
<h1>
	<img src="data:image/jpg;base64,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+a0o0PUywL+JtQIBzhYoR/7JT4PDFgsgmvGn1CcHIkvJDJj6L90flWbeaDqlh9qh0SOxmsLwsZLO6ZkEbN1KFQeD1xj6V0enQz2+m20NzIJJ441V3BJBIHPWrNFFFFFFFFFFFFVrfTrK0nlmtrSCGWY5keOMKXPuR1qzRRRRRRRRRRRRRRX/2f/hCaFodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvADw/eHBhY2tldCBiZWdpbj0n77u/JyBpZD0nVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkJz8+Cjx4OnhtcG1ldGEgeG1sbnM6eD0nYWRvYmU6bnM6bWV0YS8nPgo8cmRmOlJERiB4bWxuczpyZGY9J2h0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMnPgoKIDxyZGY6RGVzY3JpcHRpb24geG1sbnM6eG1wPSdodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvJz4KICA8eG1wOkNyZWF0b3JUb29sPkFkb2JlIFBob3Rvc2hvcCBDUzUuMSBNYWNpbnRvc2g8L3htcDpDcmVhdG9yVG9vbD4KICA8eG1wOkNyZWF0ZURhdGU+MjAxMi0wMy0wNVQxNzo1Njo1NisxMTowMDwveG1wOkNyZWF0ZURhdGU+CiAgPHhtcDpNZXRhZGF0YURhdGU+MjAxMi0wMy0wNVQxNzo1Njo1NisxMTowMDwveG1wOk1ldGFkYXRhRGF0ZT4KICA8eG1wOk1vZGlmeURhdGU+MjAxMi0wMy0wNVQxNzo1Njo1NisxMTowMDwveG1wOk1vZGlmeURhdGU+CiA8L3JkZjpEZXNjcmlwdGlvbj4KCiA8cmRmOkRlc2NyaXB0aW9uIHhtbG5zOmRjPSdodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyc+CiAgPGRjOmZvcm1hdD5pbWFnZS9qcGVnPC9kYzpmb3JtYXQ+CiA8L3JkZjpEZXNjcmlwdGlvbj4KCiA8cmRmOkRlc2NyaXB0aW9uIHhtbG5zOnhtcE1NPSdodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vJz4KICA8eG1wTU06SW5zdGFuY2VJRD54bXAuaWlkOjJFMDVBNTcxMzQyMjY4MTE4QTZERDQ3NEU5QTI0NkEyPC94bXBNTTpJbnN0YW5jZUlEPgogIDx4bXBNTTpPcmlnaW5hbERvY3VtZW50SUQ+eG1wLmRpZDoyRTA1QTU3MTM0MjI2ODExOEE2REQ0NzRFOUEyNDZBMjwveG1wTU06T3JpZ2luYWxEb2N1bWVudElEPgogIDx4bXBNTTpJbnN0YW5jZUlEPnhtcC5paWQ6MkUwNUE1NzEzNDIyNjgxMThBNkRENDc0RTlBMjQ2QTI8L3htcE1NOkluc3RhbmNlSUQ+CiAgPHhtcE1NOkRvY3VtZW50SUQgcmRmOnJlc291cmNlPSd4bXAuZGlkOjJFMDVBNTcxMzQyMjY4MTE4QTZERDQ3NEU5QTI0NkEyJyAvPgogIDx4bXBNTTpPcmlnaW5hbERvY3VtZW50SUQ+eG1wLmRpZDoyRTA1QTU3MTM0MjI2ODExOEE2REQ0NzRFOUEyNDZBMjwveG1wTU06T3JpZ2luYWxEb2N1bWVudElEPgogIDx4bXBNTTpIaXN0b3J5PgogICA8cmRmOlNlcT4KICAgPC9yZGY6U2VxPgogIDwveG1wTU06SGlzdG9yeT4KIDwvcmRmOkRlc2NyaXB0aW9uPgoKIDxyZGY6RGVzY3JpcHRpb24geG1sbnM6cGhvdG9zaG9wPSdodHRwOi8vbnMuYWRvYmUuY29tL3Bob3Rvc2hvcC8xLjAvJz4KICA8cGhvdG9zaG9wOkNvbG9yTW9kZT4xPC9waG90b3Nob3A6Q29sb3JNb2RlPgogIDxwaG90b3Nob3A6SUNDUHJvZmlsZT5Eb3QgR2FpbiAyMCU8L3Bob3Rvc2hvcDpJQ0NQcm9maWxlPgogIDxwaG90b3Nob3A6Q29sb3JNb2RlPjE8L3Bob3Rvc2hvcDpDb2xvck1vZGU+CiAgPHBob3Rvc2hvcDpJQ0NQcm9maWxlPkRvdCBHYWluIDIwJTwvcGhvdG9zaG9wOklDQ1Byb2ZpbGU+CiA8L3JkZjpEZXNjcmlwdGlvbj4KCiA8cmRmOkRlc2NyaXB0aW9uIHhtbG5zOmV4aWY9J2h0dHA6Ly9ucy5hZG9iZS5jb20vZXhpZi8xLjAvJz4KICA8ZXhpZjpPcmllbnRhdGlvbj5Ub3AtbGVmdDwvZXhpZjpPcmllbnRhdGlvbj4KICA8ZXhpZjpYUmVzb2x1dGlvbj42MDAsMDAwMDwvZXhpZjpYUmVzb2x1dGlvbj4KICA8ZXhpZjpZUmVzb2x1dGlvbj42MDAsMDAwMDwvZXhpZjpZUmVzb2x1dGlvbj4KICA8ZXhpZjpSZXNvbHV0aW9uVW5pdD5JbmNoPC9leGlmOlJlc29sdXRpb25Vbml0PgogIDxleGlmOlNvZnR3YXJlPkFkb2JlIFBob3Rvc2hvcCBDUzUuMSBNYWNpbnRvc2g8L2V4aWY6U29mdHdhcmU+CiAgPGV4aWY6RGF0ZVRpbWU+MjAxMjowMzowNSAxNzo1Njo1NjwvZXhpZjpEYXRlVGltZT4KICA8ZXhpZjpDb21wcmVzc2lvbj5KUEVHIGNvbXByZXNzaW9uPC9leGlmOkNvbXByZXNzaW9uPgogIDxleGlmOlhSZXNvbHV0aW9uPjcyPC9leGlmOlhSZXNvbHV0aW9uPgogIDxleGlmOllSZXNvbHV0aW9uPjcyPC9leGlmOllSZXNvbHV0aW9uPgogIDxleGlmOlJlc29sdXRpb25Vbml0PkluY2g8L2V4aWY6UmVzb2x1dGlvblVuaXQ+CiAgPGV4aWY6RXhpZlZlcnNpb24+RXhpZiBWZXJzaW9uIDIuMTwvZXhpZjpFeGlmVmVyc2lvbj4KICA8ZXhpZjpGbGFzaFBpeFZlcnNpb24+Rmxhc2hQaXggVmVyc2lvbiAxLjA8L2V4aWY6Rmxhc2hQaXhWZXJzaW9uPgogIDxleGlmOkNvbG9yU3BhY2U+SW50ZXJuYWwgZXJyb3IgKHVua25vd24gdmFsdWUgNjU1MzUpPC9leGlmOkNvbG9yU3BhY2U+CiAgPGV4aWY6UGl4ZWxYRGltZW5zaW9uPjEzMjY8L2V4aWY6UGl4ZWxYRGltZW5zaW9uPgogIDxleGlmOlBpeGVsWURpbWVuc2lvbj4xMDk1PC9leGlmOlBpeGVsWURpbWVuc2lvbj4KIDwvcmRmOkRlc2NyaXB0aW9uPgoKPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KPD94cGFja2V0IGVuZD0ncic/Pgr/4gOgSUNDX1BST0ZJTEUAAQEAAAOQQURCRQIQAABwcnRyR1JBWVhZWiAHzwAGAAMAAAAAAABhY3NwQVBQTAAAAABub25lAAAAAAAAAAAAAAAAAAAAAQAA9tYAAQAAAADT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">
	Callr
</h1>

<div id="ongoing" style="display: inline-block; width: 50%; vertical-align: top;" >
<div class="container">
	<h2>On Call  <span class="saving-oncall">Saving...</span></h2>
	<div id="oncall-outlet" style="min-height: 200px">
	</div>
</div>
</div
><div id="history" style="display: inline-block; width: 50%" >
<div class="container">
	<h2>People <span class="saving-people">Saving...</span></h2>
	<div id="people-outlet">
	</div
	><div id="add-people">
		<form>
		<p>
			<label>Name</label> <input name="name" type="text">
		</p>
		<p>
			<label>Email</label> <input name="email"  type="email">
		</p>
		<p>
			<label>Phone</label> <input name="phone" type="text" placeholder="+46123456789"> 
		</p>
		<button type="button" onclick="addPerson()">Add</button>

		</form>
	</div>
</div>
</div>
	


</body>
<script src="//ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
<script src="//code.jquery.com/ui/1.12.1/jquery-ui.js"></script>
<script src="//cdn.jsdelivr.net/npm/lodash@4.17.15/lodash.min.js"></script>
<script type="application/javascript">
	let PEOPLE = [];
	let ONCALL = [];
	
	function addPerson(){
	    let p = {
	        name: $('#add-people').find('input[name=name]').val(),
	        email: $('#add-people').find('input[name=email]').val(),
	        phone: $('#add-people').find('input[name=phone]').val()
	    };
	    
	    p.phone = p.phone.replace(/[^0-9^\+]/g, '');
	    if(p.phone.length < 5 || p.phone[0] !== '+'){
	        alert("Phone must be in the form +4670123456");
	        return	        
	    }
	    
	    
	    console.log(p);
	    PEOPLE.push(p);
	    savePeople()
	}
	
	function savePeople() {
	  $(".saving-people").addClass("blink");
	  setTimeout(function() {
	    $(".saving-people").removeClass("blink");
	  }, 3000)
	    
	  $.ajax({
		  type: "POST",
		  url: "/api/people",
		  data: JSON.stringify(PEOPLE),
		  contentType: "application/json"
		}).done((res) => {
	    	$("form").trigger("reset");    
	    	renderPeople();
	    })  
	}
	
	
	function renderPerson(p) {
		return '<div class="draggable person" name="' + p.phone +'" >' +
		 		'<table>' +
				'<tr>' +
					'<th>Name</th>' +
					'<td>' + p.name + '</td>' +
				'</tr>' +
				'<tr>' +
					'<th>Email</th>' +
					'<td>' + p.email + '</td>' +
				'</tr>' +
				'<tr>' +
					'<th>Phone</th>' +
					'<td>' + p.phone + '</td>' +
				'</tr>' +
				'</table>' +
				 '</div>'
	}
	
	
	function renderPeople() {
	  let dom = _.reduce(PEOPLE, (dom, p) => {
	      return dom + '<div class="p-wrapper">' +
	       '<div class="p-links">' +
	         '<a onclick="deletePerson(\'' + p.phone +'\')">delete</a> | ' +
	         '<a onclick="testCallPerson(\'' + p.phone +'\')">test call</a>' +
		   '</div>' +  
	      renderPerson(p) + 
	      '</div>'
	  }, "");
	  $("#people-outlet").html(dom);
	  mkDragDrop();
	}
	function renderOnCall() {
	  let dom = _.reduce(ONCALL, (dom, p) => dom + renderPerson(p), "");
	  
	   $("#oncall-outlet").html(dom);
		mkDragDrop();
	}
	
	function deletePerson(phone){
	    PEOPLE = _.filter(PEOPLE, (p) => p.phone !== phone);
	    ONCALL = _.filter(ONCALL, (p) => p.phone !== phone);
	    saveOncall();
	    savePeople();
	}
	
	var updateOncall = _.debounce(function() {
	    ONCALL = _.map($("#oncall-outlet").find(".person"), (o) =>{
	        console.log($(o).attr('name'));
	        return _.find(PEOPLE, {phone: $(o).attr('name')})
	    });
	  	saveOncall()
	},1000);
	
	function saveOncall() {
	  $(".saving-oncall").addClass("blink");
	  setTimeout(function() {
	    $(".saving-oncall").removeClass("blink");
	  }, 3000)
	  
	  $.ajax({
		  type: "POST",
		  url: "/api/oncall",
		  data: JSON.stringify(ONCALL),
		  contentType: "application/json"
		}).done((res) => {
	    	renderOnCall();
	    });
	}  
	
	
	
	function mkDragDrop(){
		$("#oncall-outlet").droppable({
			drop: function (event, ui) {
				 console.log('event', event);
				 console.log('ui', ui);
				 updateOncall();
			},
			out: function(event, ui ){
			    console.log("OUTSIDE");
			    console.log(event)
			    console.log(ui)
				if(event.type === 'dropout'){
				    $(ui.draggable).remove()
				    updateOncall()
				}
			},
		});
		$( "#oncall-outlet" ).sortable({
			  revert: false,	  
		});
		$( "#people-outlet .draggable" ).draggable({
		  connectToSortable: "#oncall-outlet",
		  helper: "clone",
		  revert: "invalid",
		  revertDuration: 200,
		  snap: true
		});
	}
	
	function testCallPerson(phone){
	    console.log("phone", phone)
	    $.get('api/test-call/'+phone)
	}
	
	
	function load() {
	$.get('api/people').done(people => {
	   PEOPLE = people;
	   renderPeople()
	});
	$.get('api/oncall').done(oncall => {
	   ONCALL = oncall;
	   renderOnCall() 
	});
	}
	load()
</script>

</html>
`
