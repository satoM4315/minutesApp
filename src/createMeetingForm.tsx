// https://material-ui.com/components/dialogs/#dialog
import React, { useState } from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Alert from '@material-ui/lab/Alert';

// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_CREATE_MEETING = '/meetings';

export default function CreateMeetingForm(props) {
  const [open, setOpen] = React.useState(false);

  const { onSubmitSuccessful } = props;

  const [meeting, setMeeting] = React.useState<string>();

  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const handleSubmit = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true);
    try {
      // ページが更新されないようにする
      event.preventDefault();

      // Reactのハンドラはasyncにできる
      const res = await fetch(API_URL_CREATE_MEETING, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // 相応しくないかも
        // same-originを使うべき？
        credentials: 'include',
        body: JSON.stringify({ meeting }),
      });
      const obj = await res.json();

      // 厳密だとundefinedでtrue?
      if (obj.error != null) {
        ReactDOM.render(
          <div>
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
          `An response from ${API_URL_CREATE_MEETING} returned error`
        );
      }
      // 要求は成功
      // リスナ関数を呼ぶ
      onSubmitSuccessful();
    } finally {
      setWorking(false);
      setMeeting('');
      setOpen(false);
    }
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <div>
      <Button
        variant="outlined"
        color="primary"
        disabled={working}
        onClick={handleClickOpen}
      >
        議事録の作成
      </Button>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="form-dialog-title"
      >
        <DialogTitle id="form-dialog-title">議事録作成</DialogTitle>
        <DialogContent>
          <DialogContentText>
            議事録の名前を設定してください。 注:同名の議事録は作成できません。
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="minutes name"
            type="text"
            fullWidth
            value={meeting}
            onChange={(event) => setMeeting(event.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Cancel
          </Button>
          <Button onClick={handleSubmit} color="primary">
            作成
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

CreateMeetingForm.propTypes = {
  // メッセージの更新が正常に完了したら呼ばれる関数
  onSubmitSuccessful: PropTypes.func,
};

CreateMeetingForm.defaultProps = {
  onSubmitSuccessful: () => {},
};
