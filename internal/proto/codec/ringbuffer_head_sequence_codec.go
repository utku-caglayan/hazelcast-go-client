/*
* Copyright (c) 2008-2022, Hazelcast, Inc. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License")
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package codec

import (
	"github.com/hazelcast/hazelcast-go-client/internal/proto"
)

const (
	RingbufferHeadSequenceCodecRequestMessageType  = int32(0x170300)
	RingbufferHeadSequenceCodecResponseMessageType = int32(0x170301)

	RingbufferHeadSequenceCodecRequestInitialFrameSize = proto.PartitionIDOffset + proto.IntSizeInBytes

	RingbufferHeadSequenceResponseResponseOffset = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
)

// Returns the sequence of the head. The head is the side of the ringbuffer where the oldest items in the ringbuffer
// are found. If the Ringbuffer is empty, the head will be one more than the tail.
// The initial value of the head is 0 (1 more than tail).

func EncodeRingbufferHeadSequenceRequest(name string) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(true)

	initialFrame := proto.NewFrameWith(make([]byte, RingbufferHeadSequenceCodecRequestInitialFrameSize), proto.UnfragmentedMessage)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(RingbufferHeadSequenceCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	EncodeString(clientMessage, name)

	return clientMessage
}

func DecodeRingbufferHeadSequenceResponse(clientMessage *proto.ClientMessage) int64 {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	return FixSizedTypesCodec.DecodeLong(initialFrame.Content, RingbufferHeadSequenceResponseResponseOffset)
}
