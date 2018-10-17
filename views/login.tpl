<!DOCTYPE html>
<html>
<head>
  <title>System Login</title>
  <meta http-equiv="content-type" content="text/html; charset=utf-8">
  <link rel="stylesheet" href="https://cdn.bootcss.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css"/>
</head>
<body>
<form method="post">
<div class="modal" style="display:block">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title">System Login</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        {{if .err}}
        <div class="alert alert-danger" role="alert">
          <strong>Oh snap!</strong> UserName Or Password Failed.
        </div>
        {{end}}
          <div class="form-group">
            <div class="col-md-12">
              <input type="text" name="user_name" class="form-control" placeholder="Please input user name"/>
            </div>
          </div>
          <div class="form-group">
            <div class="col-md-12">
              <input type="password"  name="password" class="form-control"  placeholder="Please input password"/>
            </div>
          </div>
      </div>
      <div class="modal-footer">
        <button type="submit" class="btn btn-primary">Submit</button>
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>
</form>
</body>
</html>
