import React, { useState } from 'react';
import PropTypes from 'prop-types';
import ReactDOM from 'react-dom';
import { createStyles, makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Alert from '@material-ui/lab/Alert';

// メッセージ追加のAPIへのURL
// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_LOGIN = '/login';

const useStyles = makeStyles((theme) =>
  createStyles({
    root: {
      '& > *': {
        margin: theme.spacing(1),
        width: '25ch',
      },
      width: '100%',
      '& > * + *': {
        marginTop: theme.spacing(2),
      },
    },
  })
);

export default function LoginPostForm(props) {
  // テキストボックス内のメッセージ
  const [userId, setUserId] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const classes = useStyles();

  const handleSubmit = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true);
    try {
      // ページが更新されないようにする
      event.preventDefault();

      // Reactのハンドラはasyncにできる
      const res = await fetch(API_URL_LOGIN, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // 相応しくないかも
        // same-originを使うべき？
        credentials: 'include',
        body: JSON.stringify({ userId, password }),
      });

      setUserId('');
      setPassword('');

      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        ReactDOM.render(
          <div className={classes.root}>
            <Alert
              variant="outlined"
              severity="error"
              onClose={() => {
                ReactDOM.render(
                  <div />,
                  document.getElementById('serverMessage')
                );
              }}
            >
              {obj.error}
            </Alert>
          </div>,
          document.getElementById('serverMessage')
        );
        throw new Error(
          `An error occurred on querying ${API_URL_LOGIN}, the response included error message: ${obj.error}`
        );
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        ReactDOM.render(
          <div className={classes.root}>
            <Alert
              variant="outlined"
              severity="error"
              onClose={() => {
                ReactDOM.render(
                  <div />,
                  document.getElementById('serverMessage')
                );
              }}
            >
              Error
            </Alert>
          </div>,
          document.getElementById('serverMessage')
        );
        throw new Error(
          `An response from ${API_URL_LOGIN} unexpectedly did not have 'success' member`
        );
      }
      if (obj.success !== true) {
        ReactDOM.render(
          <div className={classes.root}>
            <Alert
              variant="outlined"
              severity="error"
              onClose={() => {
                ReactDOM.render(
                  <div />,
                  document.getElementById('serverMessage')
                );
              }}
            >
              Error
            </Alert>
          </div>,
          document.getElementById('serverMessage')
        );
        throw new Error(
          `An response from ${API_URL_LOGIN} returned non true value as 'success' member`
        );
      }
      // 要求は成功
      ReactDOM.render(
        <div className={classes.root}>
          <Alert
            variant="outlined"
            severity="success"
            onClose={() => {
              ReactDOM.render(
                <div />,
                document.getElementById('serverMessage')
              );
            }}
          >
            ログイン成功！
          </Alert>
        </div>,
        document.getElementById('serverMessage')
      );
      // リスナ関数を呼ぶ
      props.onSubmitSuccessful();

      // ログインが成功したらmainページにリダイレクト
      window.location.href = '/';
    } finally {
      setWorking(false);
    }
  };

  return (
    <>
      <form className={classes.root} noValidate autoComplete="off">
        <TextField
          id="standard-basic"
          label="ユーザーID"
          value={userId}
          type="textbox"
          onChange={(event) => setUserId(event.target.value)}
        />
        <p />
        <TextField
          id="standard-basic"
          label="パスワード"
          value={password}
          type="password"
          onChange={(event) => setPassword(event.target.value)}
        />
        <p />
        <Button
          variant="contained"
          color="primary"
          disabled={working}
          onClick={handleSubmit}
        >
          ログイン
        </Button>
      </form>
    </>
  );
}

LoginPostForm.propTypes = {
  onSubmitSuccessful: PropTypes.func,
};

LoginPostForm.defaultProps = {
  onSubmitSuccessful: () => {},
};
